package tracing

import (
	"context"
	"os"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/pkg/errors"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"google.golang.org/grpc"
)

type Config struct {
	OtelEndpoint string `yaml:"otel_endpoint"`
	OrgID        string `yaml:"org_id"`
}

func InstallOpenTelemetryTracer(config *Config, logger log.Logger, appName, version string) (func(), error) {
	if config.OtelEndpoint == "" {
		return func() {}, nil
	}

	_ = level.Info(logger).Log("msg", "initialising OpenTelemetry tracer")

	ctx := context.Background()

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(appName),
			semconv.ServiceVersionKey.String(version),
		),
		resource.WithHost(),
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize trace resuorce")
	}

	conn, err := grpc.DialContext(ctx, config.OtelEndpoint, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, errors.Wrap(err, "failed to dial otel grpc")
	}

	options := []otlptracegrpc.Option{otlptracegrpc.WithGRPCConn(conn)}
	if config.OrgID != "" {
		options = append(options,
			otlptracegrpc.WithHeaders(map[string]string{"X-Scope-OrgID": config.OrgID}))
	}

	traceExporter, err := otlptracegrpc.New(ctx, options...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to creat trace exporter")
	}

	bsp := sdktrace.NewBatchSpanProcessor(traceExporter)
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)
	otel.SetTracerProvider(tracerProvider)

	// set global propagator to tracecontext (the default is no-op).
	otel.SetTextMapPropagator(propagation.TraceContext{})

	shutdown := func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := tracerProvider.Shutdown(ctx); err != nil {
			_ = level.Error(logger).Log("msg", "OpenTelemetry trace provider failed to shutdown", "err", err)
			os.Exit(1)
		}
	}

	return shutdown, nil
}
