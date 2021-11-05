// Code generated, do not edit
package inventory

import (
	"context"
	"os"
	"strings"

	prompt "github.com/c-bata/go-prompt"
	"github.com/olekukonko/tablewriter"
	log "github.com/sirupsen/logrus"
)

type InventoryInteractive struct {
	Inventory Inventory
}

func (i *InventoryInteractive) Executor(in string) {
	ctx := context.Background()

	in = strings.TrimSpace(in)
	blocks := strings.Split(in, " ")

	if len(blocks) < 0 {
		return
	}

	command, object, remain := blocks[0], blocks[1], blocks[2:]

	switch command {
	case "list":
		switch object {
		case "network_host":
			i.printNetworkHosts(ctx)
		case "l3_network":
			i.printL3Networks(ctx)
		case "zigbee_device":
			i.printZigbeeDevices(ctx)
		}
	case "get":
		if len(remain) < 1 {
			return
		}
		item := remain[0]

		switch object {
		case "network_host":
			i, err := i.Inventory.FetchNetworkHost(ctx, item)
			if err != nil {
				log.Error(err)
			}

			log.Infof("i: %+v", i)
		case "l3_network":
			i, err := i.Inventory.FetchL3Network(ctx, item)
			if err != nil {
				log.Error(err)
			}

			log.Infof("i: %+v", i)
		case "zigbee_device":
			i, err := i.Inventory.FetchZigbeeDevice(ctx, item)
			if err != nil {
				log.Error(err)
			}

			log.Infof("i: %+v", i)
		}
	case "create":
	case "set":
		if len(remain) < 2 {
			return
		}
		node := remain[0]
		attr := remain[1]
		val := remain[2]

		switch object {
		case "network_host":
			err := i.setNetworkHostAttribute(ctx, node, attr, val)
			if err != nil {
				log.Error(err)
			}
		case "l3_network":
			err := i.setL3NetworkAttribute(ctx, node, attr, val)
			if err != nil {
				log.Error(err)
			}
		case "zigbee_device":
			err := i.setZigbeeDeviceAttribute(ctx, node, attr, val)
			if err != nil {
				log.Error(err)
			}
		}
	}
}

func (i *InventoryInteractive) Completer(d prompt.Document) []prompt.Suggest {
	ctx := context.Background()

	blocks := strings.Split(d.CurrentLine(), " ")

	objects := []prompt.Suggest{
		{Text: "network_host", Description: "NetworkHost objects"},
		{Text: "l3_network", Description: "L3Network objects"},
		{Text: "zigbee_device", Description: "ZigbeeDevice objects"},
	}

	s := []prompt.Suggest{
		{Text: "list", Description: "List objects"},
		{Text: "get", Description: "Get an object"},
		{Text: "set", Description: "Set an object attributes"},
	}

	count := len(blocks)

	if count <= 1 {
		return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
	}

	if count == 2 {
		return prompt.FilterHasPrefix(objects, d.GetWordBeforeCursor(), true)
	}

	switch blocks[0] {
	// case "list":
	case "get":
		if count > 3 {
			return []prompt.Suggest{}
		}

		switch blocks[1] {
		case "network_host":
			return i.suggestNetworkHost(ctx, d)
		case "l3_network":
			return i.suggestL3Network(ctx, d)
		case "zigbee_device":
			return i.suggestZigbeeDevice(ctx, d)
		}
	case "set":
		if count > 4 {
			return []prompt.Suggest{}
		}

		if count == 3 {
			switch blocks[1] {
			case "network_host":
				return i.suggestNetworkHost(ctx, d)
			case "l3_network":
				return i.suggestL3Network(ctx, d)
			case "zigbee_device":
				return i.suggestZigbeeDevice(ctx, d)
			}
		}

		if count == 4 {
			switch blocks[1] {
			case "network_host":
				return i.suggestNetworkHostAttributes(ctx, d)
			case "l3_network":
				return i.suggestL3NetworkAttributes(ctx, d)
			case "zigbee_device":
				return i.suggestZigbeeDeviceAttributes(ctx, d)
			}
		}
	}

	return []prompt.Suggest{}
}
func (i *InventoryInteractive) printNetworkHosts(ctx context.Context) {
	results, err := i.Inventory.ListNetworkHosts(ctx)
	if err != nil {
		log.Error(err)
	}

	data := make([][]string, 0)

	for _, r := range results {
		data = append(data, []string{
			r.GetRole(),
			r.GetGroup(),
			r.GetName(),
			r.GetOperatingSystem(),
			r.GetPlatform(),
			r.GetType(),
			r.GetDomain(),
			r.GetDescription(),
			r.GetDn(),
		})
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{
		"role",
		"group",
		"name",
		"operating_system",
		"platform",
		"type",
		"domain",
		"description",
		"dn",
	})

	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetColumnSeparator("")

	for _, v := range data {
		table.Append(v)
	}

	table.Render()
}
func (i *InventoryInteractive) printL3Networks(ctx context.Context) {
	results, err := i.Inventory.ListL3Networks(ctx)
	if err != nil {
		log.Error(err)
	}

	data := make([][]string, 0)

	for _, r := range results {
		data = append(data, []string{
			r.GetName(),
			r.GetSoa(),
			r.GetDomain(),
			r.GetDn(),
			r.GetDescription(),
		})
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{
		"name",
		"soa",
		"domain",
		"dn",
		"description",
	})

	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetColumnSeparator("")

	for _, v := range data {
		table.Append(v)
	}

	table.Render()
}
func (i *InventoryInteractive) printZigbeeDevices(ctx context.Context) {
	results, err := i.Inventory.ListZigbeeDevices(ctx)
	if err != nil {
		log.Error(err)
	}

	data := make([][]string, 0)

	for _, r := range results {
		data = append(data, []string{
			r.GetName(),
			r.GetDescription(),
			r.GetDn(),
			r.GetIotZone(),
			r.GetType(),
			r.GetSoftwareBuildId(),
			r.GetDateCode(),
			r.GetModel(),
			r.GetVendor(),
			r.GetManufacturerName(),
			r.GetPowerSource(),
			r.GetModelId(),
		})
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{
		"name",
		"description",
		"dn",
		"iot_zone",
		"type",
		"software_build_id",
		"date_code",
		"model",
		"vendor",
		"manufacturer_name",
		"power_source",
		"model_id",
	})

	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetColumnSeparator("")

	for _, v := range data {
		table.Append(v)
	}

	table.Render()
}
func (i *InventoryInteractive) suggestNetworkHost(ctx context.Context, d prompt.Document) []prompt.Suggest {
	sugg := []prompt.Suggest{}
	results, err := i.Inventory.ListNetworkHosts(ctx)
	if err != nil {
		log.Error(err)
		return []prompt.Suggest{}
	}

	for _, r := range results {
		sugg = append(sugg, prompt.Suggest{Text: r.GetName(), Description: r.GetDescription()})
	}

	return prompt.FilterHasPrefix(sugg, d.GetWordBeforeCursor(), true)
}

func (i *InventoryInteractive) suggestNetworkHostAttributes(ctx context.Context, d prompt.Document) []prompt.Suggest {
	sugg := []prompt.Suggest{
		{Text: "role"},
		{Text: "group"},
		{Text: "name"},
		{Text: "operating_system"},
		{Text: "platform"},
		{Text: "type"},
		{Text: "domain"},
		{Text: "description"},
		{Text: "dn"},
	}
	return prompt.FilterHasPrefix(sugg, d.GetWordBeforeCursor(), true)
}

func (i *InventoryInteractive) setNetworkHostAttribute(ctx context.Context, node, attr, val string) error {
	x, err := i.Inventory.FetchNetworkHost(ctx, node)
	if err != nil {
		return err
	}

	switch attr {
	case "role":
		x.Role = val
	case "group":
		x.Group = val
	case "name":
		x.Name = val
	case "operating_system":
		x.OperatingSystem = val
	case "platform":
		x.Platform = val
	case "type":
		x.Type = val
	case "domain":
		x.Domain = val
	case "description":
		x.Description = val
	case "dn":
		x.Dn = val
	}

	x, err = i.Inventory.UpdateNetworkHost(ctx, x)
	if err != nil {
		return err
	}

	return nil
}
func (i *InventoryInteractive) suggestL3Network(ctx context.Context, d prompt.Document) []prompt.Suggest {
	sugg := []prompt.Suggest{}
	results, err := i.Inventory.ListL3Networks(ctx)
	if err != nil {
		log.Error(err)
		return []prompt.Suggest{}
	}

	for _, r := range results {
		sugg = append(sugg, prompt.Suggest{Text: r.GetName(), Description: r.GetDescription()})
	}

	return prompt.FilterHasPrefix(sugg, d.GetWordBeforeCursor(), true)
}

func (i *InventoryInteractive) suggestL3NetworkAttributes(ctx context.Context, d prompt.Document) []prompt.Suggest {
	sugg := []prompt.Suggest{
		{Text: "name"},
		{Text: "soa"},
		{Text: "domain"},
		{Text: "dn"},
		{Text: "description"},
	}
	return prompt.FilterHasPrefix(sugg, d.GetWordBeforeCursor(), true)
}

func (i *InventoryInteractive) setL3NetworkAttribute(ctx context.Context, node, attr, val string) error {
	x, err := i.Inventory.FetchL3Network(ctx, node)
	if err != nil {
		return err
	}

	switch attr {
	case "name":
		x.Name = val
	case "soa":
		x.Soa = val
	case "domain":
		x.Domain = val
	case "dn":
		x.Dn = val
	case "description":
		x.Description = val
	}

	x, err = i.Inventory.UpdateL3Network(ctx, x)
	if err != nil {
		return err
	}

	return nil
}
func (i *InventoryInteractive) suggestZigbeeDevice(ctx context.Context, d prompt.Document) []prompt.Suggest {
	sugg := []prompt.Suggest{}
	results, err := i.Inventory.ListZigbeeDevices(ctx)
	if err != nil {
		log.Error(err)
		return []prompt.Suggest{}
	}

	for _, r := range results {
		sugg = append(sugg, prompt.Suggest{Text: r.GetName(), Description: r.GetDescription()})
	}

	return prompt.FilterHasPrefix(sugg, d.GetWordBeforeCursor(), true)
}

func (i *InventoryInteractive) suggestZigbeeDeviceAttributes(ctx context.Context, d prompt.Document) []prompt.Suggest {
	sugg := []prompt.Suggest{
		{Text: "name"},
		{Text: "description"},
		{Text: "dn"},
		{Text: "iot_zone"},
		{Text: "type"},
		{Text: "software_build_id"},
		{Text: "date_code"},
		{Text: "model"},
		{Text: "vendor"},
		{Text: "manufacturer_name"},
		{Text: "power_source"},
		{Text: "model_id"},
	}
	return prompt.FilterHasPrefix(sugg, d.GetWordBeforeCursor(), true)
}

func (i *InventoryInteractive) setZigbeeDeviceAttribute(ctx context.Context, node, attr, val string) error {
	x, err := i.Inventory.FetchZigbeeDevice(ctx, node)
	if err != nil {
		return err
	}

	switch attr {
	case "name":
		x.Name = val
	case "description":
		x.Description = val
	case "dn":
		x.Dn = val
	case "iot_zone":
		x.IotZone = val
	case "type":
		x.Type = val
	case "software_build_id":
		x.SoftwareBuildId = val
	case "date_code":
		x.DateCode = val
	case "model":
		x.Model = val
	case "vendor":
		x.Vendor = val
	case "manufacturer_name":
		x.ManufacturerName = val
	case "power_source":
		x.PowerSource = val
	case "model_id":
		x.ModelId = val
	}

	x, err = i.Inventory.UpdateZigbeeDevice(ctx, x)
	if err != nil {
		return err
	}

	return nil
}
