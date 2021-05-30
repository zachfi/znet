package comms

import (
	"net/http"

	"github.com/xaque208/znet/internal/config"
)

// HTTPServerFunc is used to create a new HTTP server using a received config.
type HTTPServerFunc func(*config.Config) (*http.Server, error)

// StandardHTTPServer returns a normal HTTP server.
func StandardHTTPServer(cfg *config.Config) (*http.Server, error) {
	if cfg.HTTP == nil {
		return nil, ErrNoHTTPConfig
	}

	if cfg.HTTP.ListenAddress == "" {
		return nil, ErrNoHTTPListenAddress
	}

	return &http.Server{Addr: cfg.HTTP.ListenAddress}, nil
}
