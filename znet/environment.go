package znet

import (
	"fmt"
)

// LoadEnvironment reads reads environment variables out of vault for return.
func LoadEnvironment(config VaultConfig, env string) (map[string]string, error) {

	environment := make(map[string]string)
	if config.Host == "" || config.VaultPath == "" {
		return map[string]string{}, fmt.Errorf("Incomplete vault configuration, unable to load Environment")
	}

	// s, err := NewSecretClient(config)
	// if err != nil {
	// 	return map[string]string{}, err
	// }

	// TODO determine supported environments and their paths.
	// for _, e := range z.Config.Environments {
	// 	if e.Name == "default" {
	//
	// 		for _, k := range e.SecretValues {
	// 			path := fmt.Sprintf("%s/%s", z.Config.Vault.VaultPath, k)
	// 			log.Debugf("Reading vault path: %s", path)
	// 			secret, err := s.Logical().Read(path)
	// 			if err != nil {
	// 				log.Error(err)
	// 			}
	//
	// 			if secret != nil {
	// 				environment[k] = secret.Data["value"].(string)
	// 			}
	//
	// 		}
	//
	// 	}
	// }

	return environment, nil
}
