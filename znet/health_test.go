package znet

import (
	"io"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xaque208/znet/internal/config"
)

func TestNewServer_healthCheck(t *testing.T) {
	cfg := &config.Config{
		HTTP:  &config.HTTPConfig{},
		RPC:   &config.RPCConfig{},
		Vault: &config.VaultConfig{},
		TLS:   &config.TLSConfig{},
	}

	s, err := NewServer(cfg)
	require.NoError(t, err)
	require.NotNil(t, s)

	// Health check
	h := statusCheckHandler{server: s}

	req := httptest.NewRequest("GET", "http://example.com/foo", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)
	require.Equal(t, `{"errors":["no grpc services"],"status":"unhealthy"}`, string(body))

}
