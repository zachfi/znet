// +build unit

package lights

import (
	context "context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/xaque208/znet/internal/config"
	"github.com/xaque208/znet/pkg/iot"
)

func TestNewLights(t *testing.T) {
	l, err := NewLights(nil)
	require.Error(t, err, ErrNilConfig)
	require.Nil(t, l)

	// with config
	c := &config.Config{
		Lights: &config.LightsConfig{
			Rooms: []config.LightsRoom{},
		},
	}

	l, err = NewLights(c)
	require.NoError(t, err)
	require.NotNil(t, l)
}

func TestAddHandler(t *testing.T) {
	c := &config.Config{
		Lights: &config.LightsConfig{
			Rooms: []config.LightsRoom{},
		},
	}

	l, err := NewLights(c)
	require.NoError(t, err)
	require.NotNil(t, l)

	h := &MockLight{}

	l.AddHandler(h)

	require.Equal(t, h, l.handlers[0])
}

func TestConfiguredEventNames(t *testing.T) {

	cases := []struct {
		config *config.LightsConfig
		names  []string
		err    error
	}{
		{
			config: nil,
			names:  nil,
			err:    ErrNilConfig,
		},
		{
			config: &config.LightsConfig{
				Rooms: []config.LightsRoom{},
			},
			names: nil,
			err:   ErrNoRoomsConfigured,
		},
		{
			config: &config.LightsConfig{
				Rooms: []config.LightsRoom{
					{
						On:  []string{"one"},
						Off: []string{"two"},
					},
				},
			},
			names: []string{"one", "two"},
			err:   nil,
		},
	}

	for _, tc := range cases {
		c := &config.Config{
			Lights: tc.config,
		}

		l, err := NewLights(c)
		require.NoError(t, err)
		require.NotNil(t, l)

		names, err := l.configuredEventNames()
		require.Equal(t, tc.err, err)
		require.Equal(t, tc.names, names)

	}

}

func TestActionHandler(t *testing.T) {
	cases := map[string]struct {
		action *iot.Action
		mock   *MockLight
		err    error
		config *config.LightsConfig
	}{
		"no config": {
			config: &config.LightsConfig{},
			action: &iot.Action{},
			mock:   &MockLight{},
			err:    ErrRoomNotFound,
		},
		"simple toggle": {
			config: &config.LightsConfig{
				Rooms: []config.LightsRoom{
					{
						Name: "zone",
						On:   []string{"one"},
						Off:  []string{"two"},
					},
				},
			},
			action: &iot.Action{
				Event: "single",
				Zone:  "zone",
			},
			mock: &MockLight{
				ToggleCalls: map[string]int{"zone": 1},
			},
		},
		"double": {
			action: &iot.Action{
				Event: "double",
				Zone:  "zone1",
			},
			config: &config.LightsConfig{
				Rooms: []config.LightsRoom{
					{
						Name: "zone1",
					},
				},
			},
			mock: &MockLight{
				OnCalls:       map[string]int{"zone1": 1},
				DimCalls:      map[string]int{"zone1": 1},
				SetColorCalls: map[string]int{"zone1": 1},
			},
		},
		"triple": {
			action: &iot.Action{
				Event: "triple",
				Zone:  "zone1",
			},
			config: &config.LightsConfig{
				Rooms: []config.LightsRoom{
					{
						Name: "zone1",
					},
				},
			},
			mock: &MockLight{
				OffCalls: map[string]int{"zone1": 1},
			},
		},
		"quadruple": {
			action: &iot.Action{
				Event: "quadruple",
				Zone:  "zone1",
			},
			config: &config.LightsConfig{
				Rooms: []config.LightsRoom{
					{
						Name: "zone1",
					},
				},
			},
			mock: &MockLight{
				RandomColorCalls: map[string]int{"zone1": 1},
			},
		},
		"hold": {
			action: &iot.Action{
				Event: "hold",
				Zone:  "zone1",
			},
			config: &config.LightsConfig{
				Rooms: []config.LightsRoom{
					{
						Name: "zone1",
					},
				},
			},
			mock: &MockLight{
				DimCalls: map[string]int{"zone1": 1},
			},
		},
		"release": {
			action: &iot.Action{
				Event: "release",
				Zone:  "zone1",
			},
			config: &config.LightsConfig{
				Rooms: []config.LightsRoom{
					{
						Name: "zone1",
					},
				},
			},
			mock: &MockLight{
				DimCalls: map[string]int{"zone1": 1},
			},
		},
		"many": {
			action: &iot.Action{
				Event: "many",
				Zone:  "zone1",
			},
			config: &config.LightsConfig{
				Rooms: []config.LightsRoom{
					{
						Name: "zone1",
					},
				},
			},
			mock: &MockLight{
				AlertCalls: map[string]int{"zone1": 1},
			},
		},
	}

	for _, tc := range cases {
		h := &MockLight{}

		l := &Lights{
			config:   tc.config,
			handlers: []Handler{h},
		}

		err := l.ActionHandler(tc.action)
		require.Equal(t, tc.err, err)

		require.Equal(t, tc.mock, h)
	}

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
