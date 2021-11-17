package znet

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/xaque208/znet/internal/comms"
	"github.com/xaque208/znet/internal/config"
)

// LoadEnvironment reads reads environment variables out of vault for return.
func LoadEnvironment(cfg *config.VaultConfig, e *config.EnvironmentConfig) (map[string]string, error) {

	if e == nil {
		return map[string]string{}, fmt.Errorf("unable to load environment with nil EnvironmentConfig")
	}

	if cfg.Host == "" {
		return map[string]string{}, fmt.Errorf("a Host is required to load the environment")
	}

	if cfg.SecretRoot == "" {
		return map[string]string{}, fmt.Errorf("a SecretRoot is required to load the environment")
	}

	environment := make(map[string]string)

	s, err := comms.NewSecretClient(*cfg)
	if err != nil {
		return map[string]string{}, err
	}

	for _, k := range e.SecretValues {
		path := fmt.Sprintf("%s/%s", cfg.SecretRoot, k)
		secret, err := s.Logical().Read(path)
		if err != nil {
			log.Error(err)
		}

		if secret != nil {
			environment[k] = secret.Data["value"].(string)
		}
	}

	return environment, nil
}
