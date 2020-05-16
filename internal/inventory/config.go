package inventory

// LDAPConfig is the client configuration for LDAP.
type LDAPConfig struct {
	BaseDN string `yaml:"basedn"`
	BindDN string `yaml:"binddn"`
	BindPW string `yaml:"bindpw"`
	Host   string `yaml:"host"`
}
