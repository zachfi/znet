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

	groupName := &LightGroup{
		Name: "dungeon",
	}
	ctx := context.Background()

	l.On(ctx, groupName)
	require.Equal(t, 1, h.onCalls)

	l.Off(ctx, groupName)
	require.Equal(t, 1, h.offCalls)

	l.Alert(ctx, groupName)
	require.Equal(t, 1, h.alertCalls)

	l.Dim(ctx, groupName)
	require.Equal(t, 1, h.dimCalls)

	l.SetColor(ctx, groupName)
	require.Equal(t, 1, h.setColorCalls)

	l.Toggle(ctx, groupName)
	require.Equal(t, 1, h.toggleCalls)

	l.RandomColor(ctx, groupName)
	require.Equal(t, 1, h.randomColorCalls)
}
