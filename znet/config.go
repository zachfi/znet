package znet

import (
	"github.com/xaque208/znet/internal/agent"
	"github.com/xaque208/znet/internal/astro"
	"github.com/xaque208/znet/internal/builder"
	"github.com/xaque208/znet/internal/config"
	"github.com/xaque208/znet/internal/gitwatch"
	"github.com/xaque208/znet/internal/inventory"
	"github.com/xaque208/znet/internal/lights"
	"github.com/xaque208/znet/internal/network"
	"github.com/xaque208/znet/internal/timer"
)

// Config stores the items that are required to configure this project.
type Config struct {
	Agent        *agent.Config               `yaml:"agent,omitempty"`
	Astro        *astro.Config               `yaml:"astro,omitempty"`
	Builder      *builder.Config             `yaml:"builder,omitempty"`
	Environments *[]config.EnvironmentConfig `yaml:"environments,omitempty"`
	GitWatch     *gitwatch.Config            `yaml:"gitwatch,omitempty"`
	HTTP         *config.HTTPConfig          `yaml:"http,omitempty"`
	LDAP         *inventory.LDAPConfig       `yaml:"ldap,omitempty"`
	Lights       *lights.Config              `yaml:"lights,omitempty"`
	MQTT         *config.MQTTConfig          `yaml:"mqtt,omitempty"`
	Network      *network.Config             `yaml:"network,omitempty"`
	Rooms        *[]lights.Room              `yaml:"rooms,omitempty"`
	RPC          *config.RPCConfig           `yaml:"rpc,omitempty"`
	Timer        *timer.Config               `yaml:"timer,omitempty"`
	TLS          *config.TLSConfig           `yaml:"tls,omitempty"`
	Vault        *config.VaultConfig         `yaml:"vault,omitempty"`
}
