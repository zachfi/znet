package znet

import (
	"fmt"

	log "github.com/sirupsen/logrus"
)

func (z *Znet) LoadEnvironment() error {

	environment := make(map[string]string)

	s, err := z.NewSecretClient(z.Config.Vault)
	if err != nil {
		return err
	}

	for _, e := range z.Config.Environments {
		if e.Name == "default" {

			for _, k := range e.SecretValues {
				path := fmt.Sprintf("%s/%s", z.Config.Vault.VaultPath, k)
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
