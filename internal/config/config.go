package config

import (
	"io/ioutil"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

// LoadConfig receives a file path for a configuration to load.
func LoadConfig(file string) (*Config, error) {
	filename, _ := filepath.Abs(file)

	log.WithFields(log.Fields{
		"filename": filename,
	}).Debug("loading config file")

	config := Config{}
	err := loadYamlFile(filename, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

// loadYamlFile unmarshals a YAML file into the received interface{} or returns an error.
func loadYamlFile(filename string, d interface{}) error {
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(yamlFile, d)
	if err != nil {
		return err
	}

	return nil
}

// Config stores the items that are required to configure this project.
type Config struct {
	Agent   *AgentConfig   `yaml:"agent,omitempty"`
	Builder *BuilderConfig `yaml:"builder,omitempty"`
	// Environments *[]EnvironmentConfig `yaml:"environments,omitempty"`
	GitWatch *GitWatchConfig `yaml:"gitwatch,omitempty"`
	HTTP     *HTTPConfig     `yaml:"http,omitempty"`
	// LDAP     *LDAPConfig     `yaml:"ldap,omitempty"`
	// Lights  *LightsConfig  `yaml:"lights,omitempty"`
	// MQTT    *MQTTConfig    `yaml:"mqtt,omitempty"`
	Network *NetworkConfig `yaml:"network,omitempty"`

	TLS *TLSConfig `yaml:"tls,omitempty"`

	RPC   *RPCConfig   `yaml:"rpc,omitempty"`
	Vault *VaultConfig `yaml:"vault,omitempty"`
}

// EnvironmentConfig is the environment configuration.
type EnvironmentConfig struct {
	Name         string   `yaml:"name,omitempty"`
	SecretValues []string `yaml:"secret_values,omitempty"`
}

// RPCConfig is the configuration for the RPC client and server.
type RPCConfig struct {
	AgentListenAddress string `yaml:"agent_listen_address,omitempty"`
	ListenAddress      string `yaml:"listen_address,omitempty"`
	ServerAddress      string `yaml:"server_address,omitempty"`
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
