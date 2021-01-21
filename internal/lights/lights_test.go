package lights

import (
	context "context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xaque208/znet/internal/config"
)

func TestAlert(t *testing.T) {
	h := &mockLight{}
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
	require.Equal(t, 1, h.onCalls)

	_, err = l.Off(ctx, groupName)
	require.NoError(t, err)
	require.Equal(t, 1, h.offCalls)

	_, err = l.Alert(ctx, groupName)
	require.NoError(t, err)
	require.Equal(t, 1, h.alertCalls)

	_, err = l.Dim(ctx, groupName)
	require.NoError(t, err)
	require.Equal(t, 1, h.dimCalls)

	_, err = l.SetColor(ctx, groupName)
	require.Error(t, err)
	require.Equal(t, 0, h.setColorCalls)
	groupName.Color = "#ffffff"
	_, err = l.SetColor(ctx, groupName)
	require.NoError(t, err)
	require.Equal(t, 1, h.setColorCalls)

	_, err = l.Toggle(ctx, groupName)
	require.NoError(t, err)
	require.Equal(t, 1, h.toggleCalls)

	_, err = l.RandomColor(ctx, groupName)
	require.Error(t, err)
	require.Equal(t, 0, h.randomColorCalls)
	groupName.Colors = []string{"#ffffff"}
	_, err = l.RandomColor(ctx, groupName)
	require.NoError(t, err)
	require.Equal(t, 1, h.setColorCalls)
}
