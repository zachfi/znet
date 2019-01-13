package znet

type Data struct {
	TemplateDir   string   `yaml:"template_dir"`
	TemplatePaths []string `yaml:"template_paths"`
	DataDir       string   `yaml:"data_dir"`
	Hierarchy     []string `yaml:"hierarchy"`
}

type HostData struct {
	NTPServers            []string              `yaml:"ntp_servers" json:"ntp_servers"`
	DHCPServer            string                `yaml:"dhcp_server"`
	DHCPForwardInterfaces []string              `yaml:"dhcp_forward_interfaces"`
	LLDPInterfaces        []string              `yaml:"lldp_interfaces"`
	IRBInterfaces         []IRBInterface        `yaml:"irb_interfaces"`
	RouterAdvertisements  []RouterAdvertisement `yaml:"router_advertisements"`
	BGP                   BGP                   `yaml:"bgp"`
}

type BGP struct {
	Groups []BGPGroup `yaml:"groups"`
}

type BGPGroup struct {
	Name      string   `yaml:"name"`
	Type      string   `yaml:"type"`
	ASN       int      `yaml:"asn"`
	Neighbors []string `yaml:"neighbors"`
	Import    []string `yaml:"import"`
	Export    []string `yaml:"export"`
}

type RouterAdvertisement struct {
	Interface string `yaml:"interface"`
	DNSServer string `yaml:"dns_server"`
	Prefix    string `yaml:"prefix"`
}

type IRBInterface struct {
	Unit  string   `yaml:"unit"`
	Inet  []string `yaml:"inet"`
	Inet6 []string `yaml:"inet6"`
	MTU   int      `yaml:"mtu"`
}
