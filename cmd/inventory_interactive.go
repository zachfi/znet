// Code generated, do not edit
package cmd

import (
	"strings"

	prompt "github.com/c-bata/go-prompt"
	log "github.com/sirupsen/logrus"
	"github.com/xaque208/znet/internal/inventory"
)

type inventoryInteractive struct {
	inv inventory.Inventory
}

func (i *inventoryInteractive) executor(in string) {
	in = strings.TrimSpace(in)
	blocks := strings.Split(in, " ")

	switch blocks[0] {
	case "list":
		switch blocks[1] {
		case "network_host":
			i, err := i.inv.ListNetworkHosts()
			if err != nil {
				log.Error(err)
			}

			log.Infof("i: %+v", i)
		case "l3_network":
			i, err := i.inv.ListL3Networks()
			if err != nil {
				log.Error(err)
			}

			log.Infof("i: %+v", i)
		case "zigbee_device":
			i, err := i.inv.ListZigbeeDevices()
			if err != nil {
				log.Error(err)
			}

			log.Infof("i: %+v", i)
		case "iot_zone":
			i, err := i.inv.ListIOTZones()
			if err != nil {
				log.Error(err)
			}

			log.Infof("i: %+v", i)
		}
	case "get":
		item := blocks[2]

		switch blocks[1] {
		case "network_host":
			i, err := i.inv.FetchNetworkHost(item)
			if err != nil {
				log.Error(err)
			}

			log.Infof("i: %+v", i)
		case "l3_network":
			i, err := i.inv.FetchL3Network(item)
			if err != nil {
				log.Error(err)
			}

			log.Infof("i: %+v", i)
		case "zigbee_device":
			i, err := i.inv.FetchZigbeeDevice(item)
			if err != nil {
				log.Error(err)
			}

			log.Infof("i: %+v", i)
		case "iot_zone":
			i, err := i.inv.FetchIOTZone(item)
			if err != nil {
				log.Error(err)
			}

			log.Infof("i: %+v", i)
		}
	}

}

func (i *inventoryInteractive) completer(d prompt.Document) []prompt.Suggest {
	blocks := strings.Split(d.CurrentLine(), " ")

	objects := []prompt.Suggest{
		{Text: "network_host", Description: "NetworkHost objects"},
		{Text: "l3_network", Description: "L3Network objects"},
		{Text: "zigbee_device", Description: "ZigbeeDevice objects"},
		{Text: "iot_zone", Description: "IOTZone objects"},
	}

	s := []prompt.Suggest{
		{Text: "list", Description: "List objects"},
		{Text: "get", Description: "Get an object"},
	}

	count := len(blocks)

	if count > 0 {
		switch blocks[0] {
		case "list":
			if count > 2 {
				return []prompt.Suggest{}
			}
			return prompt.FilterHasPrefix(objects, d.GetWordBeforeCursor(), true)

		case "get":
			if count > 3 {
				return []prompt.Suggest{}
			}

			if count == 1 {
				return prompt.FilterHasPrefix(objects, d.GetWordBeforeCursor(), true)
			}

			switch blocks[1] {
			case "network_host":
				sugg := []prompt.Suggest{}
				results, err := i.inv.ListNetworkHosts()
				if err != nil {
					log.Error(err)
					return []prompt.Suggest{}
				}

				for _, r := range results {
					sugg = append(sugg, prompt.Suggest{Text: r.GetName(), Description: r.GetDescription()})
				}
				return prompt.FilterHasPrefix(sugg, d.GetWordBeforeCursor(), true)
			case "l3_network":
				sugg := []prompt.Suggest{}
				results, err := i.inv.ListL3Networks()
				if err != nil {
					log.Error(err)
					return []prompt.Suggest{}
				}

				for _, r := range results {
					sugg = append(sugg, prompt.Suggest{Text: r.GetName(), Description: r.GetDescription()})
				}
				return prompt.FilterHasPrefix(sugg, d.GetWordBeforeCursor(), true)
			case "zigbee_device":
				sugg := []prompt.Suggest{}
				results, err := i.inv.ListZigbeeDevices()
				if err != nil {
					log.Error(err)
					return []prompt.Suggest{}
				}

				for _, r := range results {
					sugg = append(sugg, prompt.Suggest{Text: r.GetName(), Description: r.GetDescription()})
				}
				return prompt.FilterHasPrefix(sugg, d.GetWordBeforeCursor(), true)
			case "iot_zone":
				sugg := []prompt.Suggest{}
				results, err := i.inv.ListIOTZones()
				if err != nil {
					log.Error(err)
					return []prompt.Suggest{}
				}

				for _, r := range results {
					sugg = append(sugg, prompt.Suggest{Text: r.GetName(), Description: r.GetDescription()})
				}
				return prompt.FilterHasPrefix(sugg, d.GetWordBeforeCursor(), true)
			}

			return prompt.FilterHasPrefix(objects, d.GetWordBeforeCursor(), true)
		}
	}

	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}
