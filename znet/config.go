package znet

import (
	"fmt"
)

// Config stores the items that are required to configure this project.
type Config struct {
	Rooms        []Room              `yaml:"rooms,omitempty"`
	Endpoint     string              `yaml:"endpoint,omitempty"`
	Environments []EnvironmentConfig `yaml:"environments,omitempty"`
	Nats         NatsConfig          `yaml:"nats,omitempty"`
	Junos        JunosConfig         `yaml:"junos,omitempty"`
	Redis        RedisConfig         `yaml:"redis,omitempty"`
	HTTP         HTTPConfig          `yaml:"http,omitempty"`
	LDAP         LDAPConfig          `yaml:"ldap,omitempty"`
	Vault        VaultConfig         `yaml:"vault,omitempty"`
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

type LDAPConfig struct {
	BaseDN string
	BindDN string
	BindPW string
	Host   string
}

type VaultConfig struct {
	Host      string
	TokenPath string `yaml:"token_path,omitempty"`
	VaultPath string `yaml:"vault_path,omitempty"`
}

type Room struct {
	Name string `yaml:"name"`
	IDs  []int  `yaml:"ids"`
}

func (c *Config) Room(name string) (Room, error) {
	for _, room := range c.Rooms {
		if room.Name == name {
			return room, nil
		}
	}

	return Room{}, fmt.Errorf("Room %s not found in config", name)
}
