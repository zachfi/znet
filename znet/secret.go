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
	var err error

	apiConfig := &api.Config{
		Address: fmt.Sprintf("https://%s:8200", config.Host),
	}

	if config.ClientKey != "" && config.ClientCert != "" {
		err = os.Setenv("VAULT_CLIENT_CERT", config.ClientCert)
		if err != nil {
			return nil, err
		}

		err = os.Setenv("VAULT_CLIENT_KEY", config.ClientKey)
		if err != nil {
			return nil, err
		}

		err = apiConfig.ReadEnvironment()
		if err != nil {
			return nil, err
		}
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
		token, err := tryCertAuth(client, config)
		if err != nil {
			log.Error(err)
		}

		if token != "" {
			client.SetToken(token)
		} else {
			return nil, fmt.Errorf("unable to summon vault token")
		}
	}

	return client, nil
}

// tryCertAuth
func tryCertAuth(client *api.Client, config VaultConfig) (string, error) {
	// https://www.vaultproject.io/api-docs/auth/cert

	log.Debugf("attempting cert authentication")
	var err error

	// to pass the password
	options := map[string]interface{}{
		"name": config.LoginName,
	}

	// the login path for cert auth
	path := "auth/cert/login"

	// PUT call to get a token
	secret, err := client.Logical().Write(path, options)
	if err != nil {
		return "", err
	}

	return secret.Auth.ClientToken, nil
}
