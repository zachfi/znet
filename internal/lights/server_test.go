// +build integration

package lights

import (
	context "context"
	"log"
	"net"
	"testing"

	"github.com/stretchr/testify/require"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"

	"github.com/xaque208/znet/internal/comms"
	"github.com/xaque208/znet/internal/config"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func TestServer(t *testing.T) {
	h := &mockLight{}

	lis = bufconn.Listen(bufSize)

	s, err := comms.TestRPCServer()
	require.NoError(t, err)
	require.NotNil(t, s)

	l := &Lights{
		config:   &config.LightsConfig{},
		handlers: []Handler{h},
	}

	RegisterLightsServer(s, l)

	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()

	defer s.Stop()

	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	defer conn.Close()
	require.NoError(t, err)
	require.NotNil(t, conn)

	client := NewLightsClient(conn)
	groupName := &LightGroupRequest{
		Name: "dungeon",
	}

	_, err = client.On(ctx, groupName)
	require.NoError(t, err)
	require.Equal(t, 1, h.onCalls)

	_, err = client.Off(ctx, groupName)
	require.NoError(t, err)
	require.Equal(t, 1, h.offCalls)

	_, err = client.Alert(ctx, groupName)
	require.NoError(t, err)
	require.Equal(t, 1, h.alertCalls)

	_, err = client.Dim(ctx, groupName)
	require.NoError(t, err)
	require.Equal(t, 1, h.dimCalls)

	_, err = client.SetColor(ctx, groupName)
	require.Error(t, err)
	require.Equal(t, 0, h.setColorCalls)
	groupName.Color = "#ffffff"
	_, err = client.SetColor(ctx, groupName)
	require.NoError(t, err)
	require.Equal(t, 1, h.setColorCalls)

	_, err = client.Toggle(ctx, groupName)
	require.NoError(t, err)
	require.Equal(t, 1, h.toggleCalls)

	_, err = client.RandomColor(ctx, groupName)
	require.Error(t, err)
	require.Equal(t, 0, h.randomColorCalls)
	groupName.Colors = []string{"#ffffff"}
	_, err = client.RandomColor(ctx, groupName)
	require.NoError(t, err)
	require.Equal(t, 1, h.setColorCalls)

}
