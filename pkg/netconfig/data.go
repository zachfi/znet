package netconfig

// Data is the structure of the data directory.
type Data struct {
	TemplateDir   string   `yaml:"template_dir"`
	TemplatePaths []string `yaml:"template_paths"`
	DataDir       string   `yaml:"data_dir"`
	Hierarchy     []string `yaml:"hierarchy"`
}

// HostData is the data relating to a particular host.
type HostData struct {
	AEInterfaces          []AEInterface         `yaml:"ae_interfaces"`
	BGP                   BGP                   `yaml:"bgp"`
	DHCPForwardInterfaces []string              `yaml:"dhcp_forward_interfaces"`
	DHCPServer            string                `yaml:"dhcp_server"`
	EthernetInterfaces    []EthernetInterface   `yaml:"eth_interfaces"`
	IRBInterfaces         []IRBInterface        `yaml:"irb_interfaces"`
	LLDPInterfaces        []string              `yaml:"lldp_interfaces"`
	NTPServers            []string              `yaml:"ntp_servers"`
	RouterAdvertisements  []RouterAdvertisement `yaml:"router_advertisements"`
	Routing               Routing               `yaml:"routing"`
	PolicyOptions         PolicyOptions         `yaml:"policy_options"`
	Security              Security              `yaml:"security"`
	VLANs                 []VLAN                `yaml:"vlans"`
}

// Security is the data related to security objects for an SRX device.
type Security struct {
	Zones          []SecurityZone         `yaml:"zones"`
	Policies       []SecurityPolicies     `yaml:"policies"`
	SimplePolicies []SimpleSecurityPolicy `yaml:"simple_policies"`
	NATRuleSets    []SecurityNATRuleSet   `yaml:"nat_rulesets"`
}

// SimpleSecurityPolicy is a simple security policy for an SRX device.
type SimpleSecurityPolicy struct {
	From string   `yaml:"from"`
	To   []string `yaml:"to"`
	Then string   `yaml:"then"`
}

// SecurityZone is a zone to which security policies are applied.
type SecurityZone struct {
	Name           string                  `yaml:"name"`
	Screen         string                  `yaml:"screen"`
	SystemServices []string                `yaml:"system_services"`
	Protocols      []string                `yaml:"protocols"`
	Interfaces     []SecurityZoneInterface `yaml:"interfaces"`
}

// SecurityPolicies is a collection of security policies.
type SecurityPolicies struct {
	From     string           `yaml:"from"`
	To       string           `yaml:"to"`
	Policies []SecurityPolicy `yaml:"policies"`
}

// SecurityPolicy is a policy to apply to a security zone.
type SecurityPolicy struct {
	Name  string   `yaml:"name"`
	Match []string `yaml:"match"`
	Then  []string `yaml:"then"`
}

// SecurityZoneInterface is a network interface to include in a security zone.
type SecurityZoneInterface struct {
	Name           string   `yaml:"name"`
	SystemServices []string `yaml:"system_services"`
	Protocols      []string `yaml:"protocols"`
}

// SecurityNATRuleSet is a set of NAT rules.
type SecurityNATRuleSet struct {
	Name  string            `yaml:"name"`
	From  string            `yaml:"from_zone"`
	To    string            `yaml:"to_zone"`
	Rules []SecurityNATRule `yaml:"rules"`
}

// SecurityNATRule is a single NAT rule.
type SecurityNATRule struct {
	Name  string               `yaml:"name"`
	Match SecurityNATRuleMatch `yaml:"match"`
}

// SecurityNATRuleMatch is a match condition for a NAT rule.
type SecurityNATRuleMatch struct {
	SourceAddressNames []string `yaml:"source_address_names"`
	SourceAddress      []string `yaml:"source_address"`
}

// BGP is the BGP stanza on a Juniper router.
type BGP struct {
	Groups []BGPGroup `yaml:"groups"`
}

// Routing is the routing configuration on a Juniper router.
type Routing struct {
	RouterID     string            `yaml:"router_id"`
	ASN          int               `yaml:"asn"`
	StaticRoutes StaticRoutes      `yaml:"static_routes"`
	Instances    []RoutingInstance `yaml:"instances"`
}

// RoutingInstance is the configuration for a Juniper routing instance.
type RoutingInstance struct {
	Name                  string       `yaml:"name"`
	Description           string       `yaml:"description"`
	InstanceType          string       `yaml:"instance_type"`
	Interfaces            []string     `yaml:"interfaces"`
	BGP                   BGP          `yaml:"bgp"`
	DHCPForwardInterfaces []string     `yaml:"dhcp_forward_interfaces"`
	DHCPServer            string       `yaml:"dhcp_server"`
	StaticRoutes          StaticRoutes `yaml:"static_routes"`
}

// PolicyOptions configures the options for a policy.
type PolicyOptions struct {
	Statements map[string]PolicyStatement `yaml:"statements"`
}

// PolicyStatement is a collection of routing policy terms.
type PolicyStatement struct {
	Name  string       `yaml:"name"`
	Terms []PolicyTerm `yaml:"terms"`
	Then  string       `yaml:"then"`
}

// PolicyTerm is a single term for a routing policy.
type PolicyTerm struct {
	From []string `yaml:"from"`
	Then string   `yaml:"then"`
}

// StaticRoutes is a collection of static routes.
type StaticRoutes struct {
	Inet  []StaticRoute `yaml:"inet"`
	Inet6 []StaticRoute `yaml:"inet6"`
}

// StaticRoute is a static route object.
type StaticRoute struct {
	Discard                    bool   `yaml:"discard"`
	NextHop                    string `yaml:"next_hop"`
	Preference                 int    `yaml:"preference"`
	Prefix                     string `yaml:"prefix"`
	QualifiedNextHopInterface  string `yaml:"qualified_next_hop_interface"`
	QualifiedNextHopPreference int    `yaml:"qualified_next_hop_preference"`
	QualifiedNextHop           string `yaml:"qualified_next_hop"`
}

// BGPGroup is a group of BGP neighbors.
type BGPGroup struct {
	Name      string   `yaml:"name"`
	Type      string   `yaml:"type"`
	ASN       int      `yaml:"asn"`
	Neighbors []string `yaml:"neighbors"`
	Import    []string `yaml:"import"`
	Export    []string `yaml:"export"`
}

// RouterAdvertisement is a router advertisement stanza for a Juniper router.
type RouterAdvertisement struct {
	Interface string `yaml:"interface"`
	DNSServer string `yaml:"dns_server"`
	Prefix    string `yaml:"prefix"`
	Managed   bool   `yaml:"managed"`
}

// IRBInterface is an Integrated Bridging and Routing interface for a Juniper router.
type IRBInterface struct {
	Unit  string   `yaml:"unit"`
	Inet  []string `yaml:"inet"`
	Inet6 []string `yaml:"inet6"`
	MTU   int      `yaml:"mtu"`
}

// InetUnit is a single interface unit for a Juniper device.
type InetUnit struct {
	Inet  []string `yaml:"inet"`
	Inet6 []string `yaml:"inet6"`
	MTU   int      `yaml:"mtu"`
}

// AEInterface is an aggregated eithernet interface on a Juniper device.
type AEInterface struct {
	Description string `yaml:"description"`
	Name        string `yaml:"name"`
	MTU         int    `yaml:"mtu"`
	Options     struct {
		MinimumLinks int      `yaml:"minimum_links"`
		LACP         []string `yaml:"lacp"`
	} `yaml:"options"`
	EthernetSwitching EthernetSwitching `yaml:"ethernet_switching"`
	Units             []InetUnit        `yaml:"units,omitempty"`
	NativeVlanID      int               `yaml:"native_vlan_id"`
}

// EthernetInterface is a single ethernet interface for a Juniper device.
type EthernetInterface struct {
	Description       string            `yaml:"description"`
	EthernetSwitching EthernetSwitching `yaml:"ethernet_switching"`
	EthernetOptions   []string          `yaml:"ethernet_options"`
	MTU               int               `yaml:"mtu"`
	Name              string            `yaml:"name"`
	NativeVlanID      int               `yaml:"native_vlan_id"`
	Units             []InetUnit        `yaml:"units"`
}

// EthernetSwitching is an ethernet switching stanza for a Juniper device.
type EthernetSwitching struct {
	Mode         string   `yaml:"mode,omitempty"`
	StormControl string   `yaml:"storm_control,omitempty"`
	VLANs        []string `yaml:"vlans,omitempty"`
}

// VLAN is a VLAN.
type VLAN struct {
	Name        string `yaml:"name"`
	ID          int    `yaml:"id"`
	Description string `yaml:"description"`
	L3Interface string `yaml:"l3_interface"`
}
