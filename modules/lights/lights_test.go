//go:build unit

package lights

import (
	"bytes"
	context "context"
	"testing"

	"github.com/go-kit/log"
	"github.com/stretchr/testify/require"

	"github.com/zachfi/znet/pkg/iot"
)

func TestNewLights(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	logger := log.NewLogfmtLogger(buf)

	c := Config{
		Rooms: []Room{},
	}

	l, err := New(c, logger)
	require.NoError(t, err)
	require.NotNil(t, l)
}

func TestAddHandler(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	logger := log.NewLogfmtLogger(buf)

	c := Config{
		Rooms: []Room{},
	}

	l, err := New(c, logger)
	require.NoError(t, err)
	require.NotNil(t, l)

	h := &MockLight{}

	l.AddHandler(h)

	require.Equal(t, h, l.handlers[0])
}

func TestConfiguredEventNames(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	logger := log.NewLogfmtLogger(buf)

	cases := []struct {
		config Config
		names  []string
		err    error
	}{
		{
			config: Config{},
			names:  nil,
			err:    ErrNoRoomsConfigured,
		},
		{
			config: Config{
				Rooms: []Room{},
			},
			names: nil,
			err:   ErrNoRoomsConfigured,
		},
		{
			config: Config{
				Rooms: []Room{
					{
						States: []StateSpec{
							{State: ZoneState_ON, Event: "one"},
							{State: ZoneState_OFF, Event: "two"},
						},
					},
				},
			},
			names: []string{"one", "two"},
			err:   nil,
		},
	}

	for _, tc := range cases {
		l, err := New(tc.config, logger)
		require.NoError(t, err)
		require.NotNil(t, l)

		names, err := l.configuredEventNames()
		require.Equal(t, tc.err, err)
		require.Equal(t, tc.names, names)
	}

}

func TestNamedTimerHandler(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	logger := log.NewLogfmtLogger(buf)

	cases := map[string]struct {
		event  string
		mock   *MockLight
		err    error
		config Config
	}{
		"no config": {
			config: Config{},
			event:  "Now",
			mock:   &MockLight{},
			err:    ErrNoRoomsConfigured,
		},
		"on": {
			event: "Later",
			config: Config{
				Rooms: []Room{
					{
						Name: "zone1",
						States: []StateSpec{
							{State: ZoneState_ON, Event: "Later"},
						},
					},
				},
			},
			mock: &MockLight{
				OnCalls:            map[string]int{"zone1": 1},
				SetBrightnessCalls: map[string]int{"zone1": 1},
				SetColorTempCalls:  map[string]int{"zone1": 1},
			},
		},
		"off": {
			event: "Later",
			config: Config{
				Rooms: []Room{
					{
						Name: "zone1",
						States: []StateSpec{
							{State: ZoneState_OFF, Event: "Later"},
						},
					},
				},
			},
			mock: &MockLight{
				OffCalls: map[string]int{"zone1": 1},
			},
		},
		"brightnes": {
			event: "Later",
			config: Config{
				Rooms: []Room{
					{
						Name: "zone1",
						States: []StateSpec{
							{State: ZoneState_ON, Event: "Later"},
						},
					},
				},
			},
			mock: &MockLight{
				OnCalls:            map[string]int{"zone1": 1},
				SetBrightnessCalls: map[string]int{"zone1": 1},
				SetColorTempCalls:  map[string]int{"zone1": 1},
			},
		},
		// "Alert": {
		// 	event: "Later",
		// 	config: Config{
		// 		Rooms: []Room{
		// 			{
		// 				Name: "zone1",
		// 				States: []StateSpec{
		// 					{State: ZoneState_ALERT, Event: "Later"},
		// 				},
		// 			},
		// 		},
		// 	},
		// 	mock: &MockLight{
		// 		AlertCalls: map[string]int{"zone1": 1},
		// 	},
		// },
		"unknown event": {
			event: "Later",
			config: Config{
				Rooms: []Room{
					{
						Name: "zone1",
					},
				},
			},
			mock: &MockLight{},
			err:  ErrUnhandledEventName,
		},
	}

	for name, tc := range cases {
		t.Logf("test: %s", name)

		h := &MockLight{}

		l, err := New(tc.config, logger)
		require.NoError(t, err)

		l.AddHandler(h)

		err = l.NamedTimerHandler(context.Background(), tc.event)
		require.Equal(t, tc.err, err)

		require.Equal(t, tc.mock, h)
	}

}

func TestActionHandler(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	logger := log.NewLogfmtLogger(buf)

	cases := map[string]struct {
		action *iot.Action
		mock   *MockLight
		err    error
		config Config
	}{
		"no config": {
			config: Config{},
			action: &iot.Action{},
			mock:   &MockLight{},
			err:    ErrUnknownActionEvent,
		},
		"simple toggle": {
			config: Config{
				Rooms: []Room{
					{
						Name: "zone",
						States: []StateSpec{
							{State: ZoneState_ON, Event: "one"},
							{State: ZoneState_OFF, Event: "two"},
						},
					},
				},
			},
			action: &iot.Action{
				Event: "single",
				Zone:  "zone",
			},
			mock: &MockLight{
				OffCalls: map[string]int{"zone": 1},
			},
		},
		"double": {
			action: &iot.Action{
				Event: "double",
				Zone:  "zone1",
			},
			config: Config{
				Rooms: []Room{
					{
						Name: "zone1",
					},
				},
			},
			mock: &MockLight{
				OnCalls:            map[string]int{"zone1": 1},
				SetBrightnessCalls: map[string]int{"zone1": 1},
				SetColorTempCalls:  map[string]int{"zone1": 1},
			},
		},
		"triple": {
			action: &iot.Action{
				Event: "triple",
				Zone:  "zone1",
			},
			config: Config{
				Rooms: []Room{
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
			config: Config{
				Rooms: []Room{
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
			config: Config{
				Rooms: []Room{
					{
						Name: "zone1",
					},
				},
			},
			mock: &MockLight{
				OnCalls:            map[string]int{"zone1": 1},
				SetBrightnessCalls: map[string]int{"zone1": 1},
				SetColorTempCalls:  map[string]int{"zone1": 1},
			},
		},
		"release": {
			action: &iot.Action{
				Event: "release",
				Zone:  "zone1",
			},
			config: Config{
				Rooms: []Room{
					{
						Name: "zone1",
					},
				},
			},
			mock: &MockLight{},
		},
		"wakeup": {
			action: &iot.Action{
				Event: "release",
				Zone:  "zone1",
			},
			config: Config{
				Rooms: []Room{
					{
						Name: "zone1",
					},
				},
			},
			mock: &MockLight{},
		},
		"many": {
			action: &iot.Action{
				Event: "many",
				Zone:  "zone1",
			},
			config: Config{
				Rooms: []Room{
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

	for name, tc := range cases {
		t.Logf("test: %s", name)
		h := &MockLight{}

		l, err := New(tc.config, logger)
		require.NoError(t, err)

		l.AddHandler(h)

		err = l.ActionHandler(context.Background(), tc.action)
		if tc.err != nil {
			require.Error(t, err, tc.err)
		}

		require.Equal(t, tc.mock, h)
	}
}
