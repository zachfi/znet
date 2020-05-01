package znet

import (
	"github.com/xaque208/znet/internal/agent"
	"github.com/xaque208/znet/internal/astro"
	"github.com/xaque208/znet/internal/builder"
	"github.com/xaque208/znet/internal/gitwatch"
	"github.com/xaque208/znet/internal/inventory"
	"github.com/xaque208/znet/internal/lights"
	"github.com/xaque208/znet/internal/timer"
)

// Config stores the items that are required to configure this project.
type Config struct {
	Agent        agent.Config         `yaml:"agent,omitempty"`
	Astro        astro.Config         `yaml:"astro,omitempty"`
	Environments []EnvironmentConfig  `yaml:"environments,omitempty"`
	GitWatch     gitwatch.Config      `yaml:"gitwatch,omitempty"`
	Builder      builder.Config       `yaml:"builder,omitempty"`
	HTTP         HTTPConfig           `yaml:"http,omitempty"`
	Junos        JunosConfig          `yaml:"junos,omitempty"`
	LDAP         inventory.LDAPConfig `yaml:"ldap,omitempty"`
	Lights       lights.Config        `yaml:"lights,omitempty"`
	Rooms        []lights.Room        `yaml:"rooms,omitempty"`
	RPC          RPCConfig            `yaml:"rpc,omitempty"`
	Timer        timer.Config         `yaml:"timer,omitempty"`
	Vault        VaultConfig          `yaml:"vault,omitempty"`
}

// EnvironmentConfig is the environment configuration.
type EnvironmentConfig struct {
	Name         string   `yaml:"name,omitempty"`
	SecretValues []string `yaml:"secret_values,omitempty"`
}

// JunosConfig is the configuration for Junos devices.
type JunosConfig struct {
	Hosts      []string
	Username   string
	PrivateKey string
}

// HTTPConfig is the configuration for the listening HTTP server.
type HTTPConfig struct {
	ListenAddress string
}

// RPCConfig is the configuration for the RPC client and server.
type RPCConfig struct {
	ListenAddress string
	ServerAddress string
}

// VaultConfig is the client configuration for Vault.
type VaultConfig struct {
	Host      string
	TokenPath string `yaml:"token_path,omitempty"`
	VaultPath string `yaml:"vault_path,omitempty"`
}
