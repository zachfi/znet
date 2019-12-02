package znet

import (
	"fmt"
)

// Config stores the items that are required to configure this project.
type Config struct {
	Rooms        []Room              `yaml:"rooms,omitempty"`
	Environments []EnvironmentConfig `yaml:"environments,omitempty"`
	Nats         NatsConfig          `yaml:"nats,omitempty"`
	Junos        JunosConfig         `yaml:"junos,omitempty"`
	Redis        RedisConfig         `yaml:"redis,omitempty"`
	HTTP         HTTPConfig          `yaml:"http,omitempty"`
	LDAP         LDAPConfig          `yaml:"ldap,omitempty"`
	Vault        VaultConfig         `yaml:"vault,omitempty"`
	RPC          RPCConfig           `yaml:"rpc,omitempty"`
	Lights       LightsConfig        `yaml:"lights,omitempty"`
}

type EnvironmentConfig struct {
	Name         string   `yaml:"name,omitempty"`
	SecretValues []string `yaml:"secret_values,omitempty"`
}

type NatsConfig struct {
	URL   string
	Topic string
}

type HueConfig struct {
	Endpoint string `yaml:"endpoint"`
	User     string `yaml:"user"`
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

type RFToyConfig struct {
	Endpoint string `yaml:"endpoint,omitempty"`
}

type LightsConfig struct {
	Rooms []Room      `yaml:"rooms"`
	Hue   HueConfig   `yaml:"hue,omitempty"`
	RFToy RFToyConfig `yaml:"rftoy,omitempty"`
}

type Room struct {
	Name   string `yaml:"name"`
	IDs    []int  `yaml:"ids"`
	HueIDs []int  `yaml:"hue"`
}

func (c *LightsConfig) Room(name string) (Room, error) {
	for _, room := range c.Rooms {
		if room.Name == name {
			return room, nil
		}
	}

	return Room{}, fmt.Errorf("Room %s not found in config", name)
}
