package znet

import (
	"io/ioutil"

	log "github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
)

func loadYamlFile(filename string, data interface{}) {
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Error(err)
	}
	err = yaml.Unmarshal(yamlFile, data)
	if err != nil {
		log.Error(err)
	}
}
