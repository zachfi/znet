// +build unit

package lights

import (
	context "context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/xaque208/znet/internal/config"
)

func TestNewLights(t *testing.T) {
	l, err := NewLights(nil)
	require.Error(t, err)
	require.Nil(t, l)
}

func TestAlert(t *testing.T) {

	testCases := []struct {
		Handler *MockLight
	}{}

	for _, tc := range testCases {
		h := &MockLight{}
		l := &Lights{
			config:   &config.LightsConfig{},
			handlers: []Handler{h},
		}

		groupName := &LightGroupRequest{
			Name: "dungeon",
		}
		ctx := context.Background()

		_, err := l.On(ctx, groupName)
		require.NoError(t, err)
		// TODO
		require.Equal(t, tc.Handler, h)

		_, err = l.Off(ctx, groupName)
		require.NoError(t, err)
		require.Equal(t, 1, h.OffCalls)

		_, err = l.Alert(ctx, groupName)
		require.NoError(t, err)
		require.Equal(t, 1, h.AlertCalls)

		_, err = l.Dim(ctx, groupName)
		require.NoError(t, err)
		require.Equal(t, 1, h.DimCalls)

		_, err = l.SetColor(ctx, groupName)
		require.Error(t, err)
		require.Equal(t, 0, h.SetColorCalls)
		groupName.Color = "#ffffff"
		_, err = l.SetColor(ctx, groupName)
		require.NoError(t, err)
		require.Equal(t, 1, h.SetColorCalls)

		_, err = l.Toggle(ctx, groupName)
		require.NoError(t, err)
		require.Equal(t, 1, h.ToggleCalls)

		_, err = l.RandomColor(ctx, groupName)
		require.Error(t, err)
		require.Equal(t, 0, h.RandomColorCalls)
		groupName.Colors = []string{"#ffffff"}
		_, err = l.RandomColor(ctx, groupName)
		require.NoError(t, err)
		require.Equal(t, 1, h.SetColorCalls)

	}

}
