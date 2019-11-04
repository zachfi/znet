package znet

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

func loadYamlFile(filename string, data interface{}) error {
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(yamlFile, data)
	if err != nil {
		return err
	}

	return nil
}
