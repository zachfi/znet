package znet

type Data struct {
	TemplateDir   string   `yaml:"template_dir"`
	TemplatePaths []string `yaml:"template_paths"`
	DataDir       string   `yaml:"data_dir"`
	Hierarchy     []string `yaml:"hierarchy"`
}

type HostData struct {
	NTPServers            []string              `yaml:"ntp_servers"`
	DHCPServer            string                `yaml:"dhcp_server"`
	DHCPForwardInterfaces []string              `yaml:"dhcp_forward_interfaces"`
	LLDPInterfaces        []string              `yaml:"lldp_interfaces"`
	IRBInterfaces         []IRBInterface        `yaml:"irb_interfaces"`
	AEInterfaces          []AEInterface         `yaml:"ae_interfaces"`
	EthernetInterfaces    []EthernetInterface   `yaml:"eth_interfaces"`
	RouterAdvertisements  []RouterAdvertisement `yaml:"router_advertisements"`
	BGP                   BGP                   `yaml:"bgp"`
	Security              Security              `yaml:"security"`
	VLANs                 []VLAN                `yaml:"vlans"`
}

type Security struct {
	Zones          []SecurityZone         `yaml:"zones"`
	Policies       []SecurityPolicies     `yaml:"policies"`
	SimplePolicies []SimpleSecurityPolicy `yaml:"simple_policies"`
	NATRuleSets    []SecurityNATRuleSet   `yaml:"nat_rulesets"`
}

type SimpleSecurityPolicy struct {
	From string   `yaml:"from"`
	To   []string `yaml:"to"`
	Then string   `yaml:"then"`
}

type SecurityZone struct {
	Name           string                  `yaml:"name"`
	Screen         string                  `yaml:"screen"`
	SystemServices []string                `yaml:"system_services"`
	Protocols      []string                `yaml:"protocols"`
	Interfaces     []SecurityZoneInterface `yaml:"interfaces"`
}

type SecurityPolicies struct {
	From     string           `yaml:"from"`
	To       string           `yaml:"to"`
	Policies []SecurityPolicy `yaml:"policies"`
}

type SecurityPolicy struct {
	Name  string   `yaml:"name"`
	Match []string `yaml:"match"`
	Then  []string `yaml:"then"`
}

type SecurityZoneInterface struct {
	Name           string   `yaml:"name"`
	SystemServices []string `yaml:"system_services"`
}

type SecurityNATRuleSet struct {
	Name  string            `yaml:"name"`
	From  string            `yaml:"from_zone"`
	To    string            `yaml:"to_zone"`
	Rules []SecurityNATRule `yaml:"rules"`
}

type SecurityNATRule struct {
	Name  string               `yaml:"name"`
	Match SecurityNATRuleMatch `yaml:"match"`
}

type SecurityNATRuleMatch struct {
	SourceAddressNames []string `yaml:"source_address_names"`
	SourceAddress      []string `yaml:"source_address"`
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

type AEInterface struct {
	Description string `yaml:"description"`
	Name        string `yaml:"name"`
	MTU         int    `yaml:"mtu"`
	Options     struct {
		MinimumLinks int      `yaml:"minimum_links"`
		LACP         []string `yaml:"lacp"`
	} `yaml:"options"`
	EthernetSwitching EthernetSwitching `yaml:"ethernet_switching"`
}

type EthernetInterface struct {
	Description       string            `yaml:"description"`
	EthernetSwitching EthernetSwitching `yaml:"ethernet_switching"`
	EthernetOptions   []string          `yaml:"ethernet_options"`
	MTU               int               `yaml:"mtu"`
	Name              string            `yaml:"name"`
	NativeVlanId      int               `yaml:"native_vlan_id"`
}

type EthernetSwitching struct {
	Mode         string   `yaml:"mode"`
	StormControl string   `yaml:"storm_control"`
	VLANs        []string `yaml:"vlans"`
}

type VLAN struct {
	Name        string `yaml:"name"`
	ID          int    `yaml:"id"`
	Description string `yaml:"description"`
	L3Interface string `yaml:"l3_interface"`
}
