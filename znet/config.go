package znet

import (
	"github.com/xaque208/znet/internal/lights"
	"github.com/xaque208/znet/internal/timer"
)

// Config stores the items that are required to configure this project.
type Config struct {
	Rooms        []lights.Room       `yaml:"rooms,omitempty"`
	Environments []EnvironmentConfig `yaml:"environments,omitempty"`
	Nats         NatsConfig          `yaml:"nats,omitempty"`
	Junos        JunosConfig         `yaml:"junos,omitempty"`
	Redis        RedisConfig         `yaml:"redis,omitempty"`
	HTTP         HTTPConfig          `yaml:"http,omitempty"`
	LDAP         LDAPConfig          `yaml:"ldap,omitempty"`
	Vault        VaultConfig         `yaml:"vault,omitempty"`
	RPC          RPCConfig           `yaml:"rpc,omitempty"`
	Lights       lights.LightsConfig `yaml:"lights,omitempty"`
	Events       EventsConfig        `yaml:"events,omitempty"`
	Timer        timer.TimerConfig   `yaml:"timer,omitempty"`
}

type EnvironmentConfig struct {
	Name         string   `yaml:"name,omitempty"`
	SecretValues []string `yaml:"secret_values,omitempty"`
}

type NatsConfig struct {
	URL   string
	Topic string
}

type RedisConfig struct {
	Host string
}

type JunosConfig struct {
	Hosts      []string
	Username   string
	PrivateKey string
}

type HTTPConfig struct {
	ListenAddress string
}

type RPCConfig struct {
	ListenAddress string
	ServerAddress string
}

type LDAPConfig struct {
	BaseDN    string `yaml:"basedn"`
	BindDN    string `yaml:"binddn"`
	BindPW    string `yaml:"bindpw"`
	Host      string `yaml:"host"`
	UnknownDN string `yaml:"unknowndn"`
}

type VaultConfig struct {
	Host      string
	TokenPath string `yaml:"token_path,omitempty"`
	VaultPath string `yaml:"vault_path,omitempty"`
}

type EventsConfig struct {
}
