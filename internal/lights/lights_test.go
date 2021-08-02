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
	c := &config.LightsConfig{
		Rooms: []config.LightsRoom{},
	}

	l, err = NewLights(c)
	require.NoError(t, err)
	require.NotNil(t, l)
}

func TestAddHandler(t *testing.T) {
	c := &config.LightsConfig{
		Rooms: []config.LightsRoom{},
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
		l, err := NewLights(tc.config)
		if tc.config == nil {
			require.Error(t, err, ErrNilConfig)
			require.Nil(t, l)
			continue
		} else {
			require.NoError(t, err)
			require.NotNil(t, l)
		}

		names, err := l.configuredEventNames()
		require.Equal(t, tc.err, err)
		require.Equal(t, tc.names, names)

	}

}

func TestNamedTimerHandler(t *testing.T) {
	cases := map[string]struct {
		event  string
		mock   *MockLight
		err    error
		config *config.LightsConfig
	}{
		"no config": {
			config: &config.LightsConfig{},
			event:  "Now",
			mock:   &MockLight{},
			err:    ErrNoRoomsConfigured,
		},
		"On": {
			event: "Later",
			config: &config.LightsConfig{
				Rooms: []config.LightsRoom{
					{
						Name: "zone1",
						On:   []string{"Later"},
					},
				},
			},
			mock: &MockLight{
				OnCalls: map[string]int{"zone1": 1},
			},
		},
		"Off": {
			event: "Later",
			config: &config.LightsConfig{
				Rooms: []config.LightsRoom{
					{
						Name: "zone1",
						Off:  []string{"Later"},
					},
				},
			},
			mock: &MockLight{
				OffCalls: map[string]int{"zone1": 1},
			},
		},
		"Dim": {
			event: "Later",
			config: &config.LightsConfig{
				Rooms: []config.LightsRoom{
					{
						Name: "zone1",
						Dim:  []string{"Later"},
					},
				},
			},
			mock: &MockLight{
				DimCalls: map[string]int{"zone1": 1},
			},
		},
		"Alert": {
			event: "Later",
			config: &config.LightsConfig{
				Rooms: []config.LightsRoom{
					{
						Name:  "zone1",
						Alert: []string{"Later"},
					},
				},
			},
			mock: &MockLight{
				AlertCalls: map[string]int{"zone1": 1},
			},
		},
		"unknown event": {
			event: "Later",
			config: &config.LightsConfig{
				Rooms: []config.LightsRoom{
					{
						Name: "zone1",
					},
				},
			},
			mock: &MockLight{},
			err:  ErrUnhandledEventName,
		},
	}

	for _, tc := range cases {
		h := &MockLight{}

		l, err := NewLights(tc.config)
		require.NoError(t, err)

		l.AddHandler(h)

		err = l.NamedTimerHandler(context.Background(), tc.event)
		require.Equal(t, tc.err, err)

		require.Equal(t, tc.mock, h)
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

		l, err := NewLights(tc.config)
		require.NoError(t, err)

		l.AddHandler(h)

		err = l.ActionHandler(tc.action)
		if tc.err != nil {
			require.Contains(t, err.Error(), tc.err.Error())
		}

		require.Equal(t, tc.mock, h)
	}

}

func TestAlert(t *testing.T) {

	testCases := []struct {
		Handler *MockLight
	}{}

	for _, tc := range testCases {
		h := &MockLight{}
		l, err := NewLights(&config.LightsConfig{})
		require.NoError(t, err)

		l.AddHandler(h)

		groupName := &LightGroupRequest{
			Name: "dungeon",
		}
		ctx := context.Background()

		_, err = l.On(ctx, groupName)
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
