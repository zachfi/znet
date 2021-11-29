package inventory

// Config is the client configuration for LDAP.
type Config struct {
	LDAP LDAPConfig `yaml:"ldap,omitempty"`
}

type LDAPConfig struct {
	BaseDN string `yaml:"basedn"`
	BindDN string `yaml:"binddn"`
	BindPW string `yaml:"bindpw"`
	Host   string `yaml:"host"`
}
