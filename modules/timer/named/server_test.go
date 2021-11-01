package named

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/xaque208/znet/internal/config"
	"github.com/xaque208/znet/modules/lights"
)

func TestServer(t *testing.T) {

	c := &config.LightsConfig{
		Rooms: []config.LightsRoom{
			{
				On: []string{"one"},
			},
		},
	}

	l, err := lights.NewLights(c)
	require.NoError(t, err)
	require.NotNil(t, l)

	s, err := NewServer(l)
	require.NoError(t, err)
	require.NotNil(t, s)

	// NamedTimer
	req := &NamedTimeStamp{}
	e, err := s.NamedTimer(context.Background(), req)
	require.Error(t, err)
	require.Nil(t, e)

	req = &NamedTimeStamp{
		Name: "sunrise",
	}
	e, err = s.NamedTimer(context.Background(), req)
	require.Equal(t, lights.ErrUnhandledEventName, err)
	require.NotNil(t, e)
}
