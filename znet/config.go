package znet

import (
	"fmt"
)

type Config struct {
	Rooms    []Room      `yaml:"rooms"`
	Endpoint string      `yaml:"endpoint"`
	Nats     NatsConfig  `yaml:"nats,omitempty"`
	Junos    NatsConfig  `yaml:"junos,omitempty"`
	Redis    RedisConfig `yaml:"redis,omitempty"`
	Http     HttpConfig  `yaml:"http,omitempty"`
	Ldap     LdapConfig  `yaml:"ldap,omitempty"`
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

type HttpConfig struct {
	ListenAddress string
}

type LdapConfig struct {
	BaseDN string
	BindDN string
	BindPW string
	Host   string
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
