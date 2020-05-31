package znet

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/hashicorp/vault/api"
	log "github.com/sirupsen/logrus"
)

// NewSecretClient receives a configuration and returns a client for Vault.
func NewSecretClient(config VaultConfig) (*api.Client, error) {
	apiConfig := &api.Config{
		Address: fmt.Sprintf("https://%s:8200", config.Host),
	}

	client, err := api.NewClient(apiConfig)
	if err != nil {
		return &api.Client{}, err
	}

	envToken := os.Getenv("VAULT_TOKEN")

	if envToken != "" {
		client.SetToken(envToken)
	} else if config.TokenPath != "" {
		token, err := ioutil.ReadFile(config.TokenPath)
		if err != nil {
			log.Error(err)
		}

		client.SetToken(string(token))
	} else {
		return nil, fmt.Errorf("no vault client token set")
	}

	return client, nil
}
