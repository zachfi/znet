package inventory

import (
	"testing"

	prompt "github.com/c-bata/go-prompt"
	"github.com/stretchr/testify/require"
)

func TestCompleter_Completer(t *testing.T) {
	inv := &MockInventory{}

	i := InventoryInteractive{
		Inventory: inv,
	}

	inv.ListZigbeeDeviceResponse = []ZigbeeDevice{
		{
			Name: "0x0",
		},
	}

	commands := []prompt.Suggest{
		{Text: "list", Description: "List objects"},
		{Text: "get", Description: "Get an object"},
		{Text: "set", Description: "Set an object attributes"},
	}

	objects := []prompt.Suggest{
		{
			Text:        "network_host",
			Description: "NetworkHost objects",
		},
		{
			Text:        "l3_network",
			Description: "L3Network objects",
		},
		{
			Text:        "zigbee_device",
			Description: "ZigbeeDevice objects",
		},
	}

	zigbeeDevices := []prompt.Suggest{
		{
			Text: "0x0",
		},
	}

	zigbeeAttrs := []prompt.Suggest{
		{Text: "name", Description: ""},
		{Text: "description", Description: ""},
		{Text: "dn", Description: ""},
		{Text: "iot_zone", Description: ""},
		{Text: "type", Description: ""},
		{Text: "software_build_id", Description: ""},
		{Text: "date_code", Description: ""},
		{Text: "model", Description: ""},
		{Text: "vendor", Description: ""},
		{Text: "manufacturer_name", Description: ""},
		{Text: "power_source", Description: ""},
		{Text: "model_id", Description: ""}}

	cases := []struct {
		document prompt.Document
		results  []prompt.Suggest
	}{
		{
			document: prompt.Document{Text: ""},
			results:  commands,
		},
		{
			document: prompt.Document{Text: "list "},
			results:  objects,
		},
		{
			document: prompt.Document{Text: "get "},
			results:  objects,
		},
		{
			document: prompt.Document{Text: "get zigbee_device "},
			results:  zigbeeDevices,
		},
		{
			document: prompt.Document{Text: "get zigbee_device x "},
			results:  []prompt.Suggest{},
		},
		{
			document: prompt.Document{Text: "set "},
			results:  objects,
		},
		{
			document: prompt.Document{Text: "set zigbee_device "},
			results:  zigbeeDevices,
		},
		{
			document: prompt.Document{Text: "set zigbee_device 0x0 "},
			results:  zigbeeAttrs,
		},
	}

	for _, tc := range cases {
		sugg := i.Completer(tc.document)
		require.Equal(t, tc.results, sugg)
	}
}

func TestCompleter_Executor(t *testing.T) {
	inv := &MockInventory{}

	i := &InventoryInteractive{
		Inventory: inv,
	}

	inv.ListZigbeeDeviceResponse = []ZigbeeDevice{
		{
			Name: "0x0",
		},
	}

	inv.FetchZigbeeDeviceResponse = &ZigbeeDevice{
		Name: "0x0",
	}

	cases := []struct {
		text string
		mock MockInventory
	}{
		{
			text: "set zigbee_device 0x0 vendor test",
			mock: MockInventory{
				FetchZigbeeDeviceCalls:  map[string]int{"0x0": 1},
				UpdateZigbeeDeviceCalls: map[string]int{"0x0": 1},
			},
		},
	}

	// FetchZigbeeDeviceCalls    map[string]int
	// FetchZigbeeDeviceResponse *ZigbeeDevice
	// FetchZigbeeDeviceErr      error
	// ListZigbeeDeviceResponse  []ZigbeeDevice
	// ListZigbeeDeviceErr       error
	// CreateZigbeeDeviceCalls   map[string]int

	for _, tc := range cases {
		i.Executor(tc.text)
		require.Equal(t, tc.mock.FetchZigbeeDeviceCalls, inv.FetchZigbeeDeviceCalls)
		require.Equal(t, tc.mock.UpdateZigbeeDeviceCalls, inv.UpdateZigbeeDeviceCalls)
	}

}
