package comms

import (
	"crypto/tls"
	"strings"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/xaque208/znet/internal/config"
)

// StandardRPCClient implements enough to get standard gRPC client connection.
func StandardRPCClient(serverAddress string, cfg config.Config, logger log.Logger) *grpc.ClientConn {
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

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)),
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
		grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor()),
	}

	conn, err := grpc.Dial(serverAddress, opts...)
	if err != nil {
		_ = level.Error(logger).Log("msg", "failed dialing gRPC", "err", err)
	}

	return conn
}

func SlimRPCClient(serverAddress string, logger log.Logger) *grpc.ClientConn {
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
		grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor()),
	}

	_ = level.Debug(logger).Log("msg", "dialing gRPC", "server_address", serverAddress)

	conn, err := grpc.Dial(serverAddress, opts...)
	if err != nil {
		_ = level.Error(logger).Log("msg", "failed dialing gRPC", "err", err)
	}

	return conn
}
