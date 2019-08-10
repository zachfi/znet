package znet

import (
	"io/ioutil"

	"github.com/hashicorp/vault/api"
	log "github.com/sirupsen/logrus"
)

// NewSecretClient receives a configuration and returns a client for Vault.
func (z *Znet) NewSecretClient(config VaultConfig) (*api.Client, error) {

	apiConfig := &api.Config{
		Address: config.Host,
	}

	client, err := api.NewClient(apiConfig)
	if err != nil {
		return &api.Client{}, err
	}

	token, err := ioutil.ReadFile(z.Config.Vault.TokenPath)
	if err != nil {
		log.Error(err)
	}

	client.SetToken(string(token))

	return client, nil
}
