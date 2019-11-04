package znet

import (
	"fmt"

	log "github.com/sirupsen/logrus"
)

func (z *Znet) LoadEnvironment() error {

	environment := make(map[string]string)
	if z.Config.Vault.Host == "" || z.Config.Vault.VaultPath == "" {
		return fmt.Errorf("Incomplete vault configuration, unable to load Environment")
	}

	s, err := z.NewSecretClient(z.Config.Vault)
	if err != nil {
		return err
	}

	for _, e := range z.Config.Environments {
		if e.Name == "default" {

			for _, k := range e.SecretValues {
				path := fmt.Sprintf("%s/%s", z.Config.Vault.VaultPath, k)
				log.Debugf("Reading vault path: %s", path)
				secret, err := s.Logical().Read(path)
				if err != nil {
					log.Error(err)
				}

				if secret != nil {
					environment[k] = secret.Data["value"].(string)
				}

			}

		}
	}

	z.Environment = environment

	return nil
}
