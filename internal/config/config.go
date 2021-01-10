package config

// Config stores the items that are required to configure this project.
type Config struct {
	// Agent        *agent.Config         `yaml:"agent,omitempty"`
	// Astro        *astro.Config         `yaml:"astro,omitempty"`
	// Builder      *builder.Config       `yaml:"builder,omitempty"`
	Environments *[]EnvironmentConfig `yaml:"environments,omitempty"`
	// GitWatch     *gitwatch.Config     `yaml:"gitwatch,omitempty"`
	HTTP *HTTPConfig `yaml:"http,omitempty"`
	// LDAP         *inventory.LDAPConfig `yaml:"ldap,omitempty"`
	// Lights       *lights.Config        `yaml:"lights,omitempty"`
	MQTT *MQTTConfig `yaml:"mqtt,omitempty"`
	// Network      *network.Config       `yaml:"network,omitempty"`
	// Rooms        *[]lights.Room        `yaml:"rooms,omitempty"`
	RPC *RPCConfig `yaml:"rpc,omitempty"`
	// Timer        *timer.Config         `yaml:"timer,omitempty"`
	TLS   *TLSConfig   `yaml:"tls,omitempty"`
	Vault *VaultConfig `yaml:"vault,omitempty"`
}

// EnvironmentConfig is the environment configuration.
type EnvironmentConfig struct {
	Name         string   `yaml:"name,omitempty"`
	SecretValues []string `yaml:"secret_values,omitempty"`
}

// RPCConfig is the configuration for the RPC client and server.
type RPCConfig struct {
	ListenAddress string `yaml:"listen_address,omitempty"`
	ServerAddress string `yaml:"server_address,omitempty"`
}

// HTTPConfig is the configuration for the listening HTTP server.
type HTTPConfig struct {
	ListenAddress string `yaml:"listen_address,omitempty"`
}

// TLSConfig is the configuration for an RPC TLS client and server.
type TLSConfig struct {
	// CN is the common name to use when issuing a new certificate.
	CN string `yaml:"cn"`

	// CAFile is the file path of the CA for vault HTTPs certificate.
	CAFile string `yaml:"ca_file"`

	// CacheDir is the directory to cache the TLS files in.
	CacheDir string `yaml:"cache_dir"`
}

// VaultConfig is the client configuration for Vault.
type VaultConfig struct {
	Host       string `yaml:"host,omitempty"`
	TokenPath  string `yaml:"token_path,omitempty"`
	SecretRoot string `yaml:"secret_root,omitempty"`

	ClientKey  string `yaml:"client_key,omitempty"`
	ClientCert string `yaml:"client_cert,omitempty"`
	CACert     string `yaml:"ca_cert,omitempty"`
	LoginName  string `yaml:"login_name,omitempty"`
}

type MQTTConfig struct {
	URL      string `yaml:"url,omitempty"`
	Topic    string `yaml:"topic,omitempty"`
	Username string `yaml:"username,omitempty"`
	Password string `yaml:"password,omitempty"`
}
