//go:build integration

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
	h := &MockLight{}

	lis = bufconn.Listen(bufSize)

	s, err := comms.TestRPCServer()
	require.NoError(t, err)
	require.NotNil(t, s)

	l, err := NewLights(&config.LightsConfig{})
	require.NoError(t, err)

	l.AddHandler(h)

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

	testCases := []struct {
		Call    func(context.Context, *LightGroupRequest, ...grpc.CallOption) (*LightResponse, error)
		Handler *MockLight
	}{
		{
			Call: client.On,
			Handler: &MockLight{
				OnCalls: map[string]int{"dungeon": 1},
			},
		},
		{
			Call: client.Off,
			Handler: &MockLight{
				// TODO fix this cary over from the previous itteration of the loop so we don't have to check the accumulation of the results here.
				OnCalls:  map[string]int{"dungeon": 1},
				OffCalls: map[string]int{"dungeon": 1},
			},
		},
	}

	for _, tc := range testCases {
		_, err = tc.Call(ctx, groupName)
		require.NoError(t, err)
		require.Equal(t, tc.Handler, h)

		// _, err = client.Alert(ctx, groupName)
		// require.NoError(t, err)
		// require.Equal(t, 1, h.AlertCalls)
		//
		// _, err = client.Dim(ctx, groupName)
		// require.NoError(t, err)
		// require.Equal(t, 1, h.DimCalls)
		//
		// _, err = client.SetColor(ctx, groupName)
		// require.Error(t, err)
		// require.Equal(t, 0, h.SetColorCalls)
		// groupName.Color = "#ffffff"
		// _, err = client.SetColor(ctx, groupName)
		// require.NoError(t, err)
		// require.Equal(t, 1, h.SetColorCalls)
		//
		// _, err = client.Toggle(ctx, groupName)
		// require.NoError(t, err)
		// require.Equal(t, 1, h.ToggleCalls)
		//
		// _, err = client.RandomColor(ctx, groupName)
		// require.Error(t, err)
		// require.Equal(t, 0, h.RandomColorCalls)
		// groupName.Colors = []string{"#ffffff"}
		// _, err = client.RandomColor(ctx, groupName)
		// require.NoError(t, err)
		// require.Equal(t, 1, h.SetColorCalls)
	}

}
