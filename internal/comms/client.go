package comms

import (
	"crypto/tls"
	"strings"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	otgrpc "github.com/opentracing-contrib/go-grpc"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/xaque208/znet/internal/config"
)

// StandardRPCClient implements enough to get standard gRPC client connection.
func StandardRPCClient(serverAddress string, cfg config.Config, logger log.Logger) *grpc.ClientConn {
	var opts []grpc.DialOption

	roots, err := CABundle(cfg.Vault)
	if err != nil {
		_ = level.Error(logger).Log("err", err)
	}

	c, err := newCertify(cfg.Vault, cfg.TLS)
	if err != nil {
		_ = level.Error(logger).Log("err", err)
	}

	serverName := strings.Split(serverAddress, ":")[0]

	tlsConfig := &tls.Config{
		ServerName:           serverName,
		InsecureSkipVerify:   false,
		RootCAs:              roots,
		GetClientCertificate: c.GetClientCertificate,
		GetCertificate:       c.GetCertificate,
	}

	// opts = append(opts, grpc.WithInsecure())
	opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)))
	// opts = append(opts, grpc.WithBlock())
	opts = append(opts, grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())))
	opts = append(opts, grpc.WithStreamInterceptor(otgrpc.OpenTracingStreamClientInterceptor(opentracing.GlobalTracer())))

	conn, err := grpc.Dial(serverAddress, opts...)
	if err != nil {
		_ = level.Error(logger).Log("msg", "failed dialing gRPC", "err", err)
	}

	return conn
}

func SlimRPCClient(serverAddress string, logger log.Logger) *grpc.ClientConn {
	var opts []grpc.DialOption

	opts = append(opts, grpc.WithInsecure())
	// opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)))
	// opts = append(opts, grpc.WithBlock())
	opts = append(opts, grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())))
	opts = append(opts, grpc.WithStreamInterceptor(otgrpc.OpenTracingStreamClientInterceptor(opentracing.GlobalTracer())))

	_ = level.Debug(logger).Log("msg", "dialing gRPC", "server_address", serverAddress)

	conn, err := grpc.Dial(serverAddress, opts...)
	if err != nil {
		_ = level.Error(logger).Log("msg", "failed dialing gRPC", "err", err)
	}

	return conn
}
