package comms

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/hashicorp/vault/api"
	log "github.com/sirupsen/logrus"

	"github.com/xaque208/znet/internal/config"
)

// NewSecretClient receives a configuration and returns a client for Vault.
func NewSecretClient(config config.VaultConfig) (*api.Client, error) {
	var err error
	var token string

	if config.Host == "" {
		return nil, ErrMissingVaultConfig
	}

	apiConfig := &api.Config{
		Address: config.Host,
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
		token = envToken
	} else if config.TokenPath != "" {
		cachedToken, err := ioutil.ReadFile(config.TokenPath)
		if err != nil {
			log.Error(err)
		}

		token = string(cachedToken)

		// Once we've loaded a token from the file, we should validate it before
		// moving on.  This gives us an opportinutiy to replace the token with a
		// valid one using the cert auth below.  If we receive an error here, we
		// clear the token to allow us to proceed.
		err = validateToken(client, token)
		if err != nil {
			log.Error(err)
			token = ""
			client.ClearToken()
		}
	}

	if token == "" {
		certAuthToken, err := tryCertAuth(client, config)
		if err != nil {
			log.Error(err)
		}

		if certAuthToken != "" {
			token = certAuthToken

			if config.TokenPath != "" {
				err = saveToken(token, config.TokenPath)
				if err != nil {
					log.Error(err)
				}
			}
		}
	}

	if token != "" {
		client.SetToken(token)
	} else {
		return nil, fmt.Errorf("unable to summon vault token")
	}

	return client, nil
}

// tryCertAuth
func tryCertAuth(client *api.Client, config config.VaultConfig) (string, error) {
	// https://www.vaultproject.io/api-docs/auth/cert

	log.Debug("attempting cert authentication")
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

func saveToken(token string, tokenPath string) error {
	log.WithFields(log.Fields{
		"token_path": tokenPath,
	}).Debugf("saving token")

	f, err := os.Create(tokenPath)
	if err != nil {
		return err
	}
	defer f.Close()

	err = f.Chmod(0600)
	if err != nil {
		return err
	}

	_, err = f.WriteString(token)
	if err != nil {
		return err
	}

	return nil
}

// validateToken ensures that a token is still valid.  The caller is resposible
// for setting the token on the client again.  Any errors are returned.
func validateToken(client *api.Client, token string) error {
	var err error

	client.SetToken(token)
	defer client.ClearToken()

	t := client.Auth().Token()
	if err != nil {
		return err
	}

	s, err := t.LookupSelf()
	if err != nil {
		return err
	}

	var expireTime time.Time
	if s != nil {
		if expireTimeStamp, ok := s.Data["expire_time"]; ok {
			expireTime, err = time.Parse(time.RFC3339, expireTimeStamp.(string))
			if err != nil {
				return err
			}
		}

		log.WithFields(log.Fields{
			"expire_time": expireTime.Format(time.RFC3339),
		}).Trace("vault token")

		if time.Until(expireTime) < 60*time.Second {
			return fmt.Errorf("token expires soon: %s", expireTime.Format(time.RFC3339))
		}
	}

	return nil
}
