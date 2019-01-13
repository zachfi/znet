package znet

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/alecthomas/template"
	log "github.com/sirupsen/logrus"
)

type Znet struct {
	ConfigDir string
	Config    Config
	Data      Data
}

func (z *Znet) LoadConfig(file string) {
	filename, _ := filepath.Abs(file)
	log.Debugf("Loading config from: %s", filename)
	config := Config{}
	loadYamlFile(filename, &config)

	z.Config = config
}

func (z *Znet) LoadData(configDir string) {
	log.Debugf("Loading data from: %s", configDir)
	dataConfig := Data{}
	loadYamlFile(fmt.Sprintf("%s/%s", configDir, "data.yaml"), &dataConfig)

	z.Data = dataConfig
}

func (z *Znet) ConfigureNetworkHost(host *NetworkHost, commit bool) {
	// log.Warnf("Znet: %+v", z)
	// log.Warnf("Commit: %t", commit)
	// log.Warnf("Host: %+v", host)
	templates := z.TemplatesForDevice(*host)
	log.Infof("Templates for host %s: %+v", host.Name, templates)

	for _, t := range templates {
		result := z.RenderHostTemplate(*host, t)
		log.Infof("Result: %+v", result)

	}
	// log.Infof("Host: %+v", host)
}

func (z *Znet) TemplatePathsForDevice(host NetworkHost) []string {
	var templates []string

	for _, t := range z.Data.TemplatePaths {
		tmpl, err := template.New("test").Parse(t)

		var buf bytes.Buffer

		err = tmpl.Execute(&buf, host)
		if err != nil {
			log.Error(err)
		}

		templates = append(templates, buf.String())
	}

	return templates
}

func (z *Znet) TemplatesForDevice(host NetworkHost) []string {
	var files []string

	paths := z.TemplatePathsForDevice(host)

	for _, p := range paths {
		templateAbs := fmt.Sprintf("%s/%s", z.ConfigDir, p)
		if _, err := os.Stat(templateAbs); err == nil {
			globPattern := fmt.Sprintf("%s/*.tmpl", templateAbs)
			foundFiles, err := filepath.Glob(globPattern)
			if err != nil {
				log.Error(err)
			} else {
				for _, f := range foundFiles {
					files = append(files, f)
				}
			}

		} else if os.IsNotExist(err) {
			log.Warnf("Template path %s does not exist", templateAbs)
		}

	}

	return files
}

func (z *Znet) RenderHostTemplate(host NetworkHost, path string) string {

	b, err := ioutil.ReadFile(path)
	if err != nil {
		log.Error(err)
	}

	str := string(b)
	tmpl, err := template.New("test").Parse(str)

	var buf bytes.Buffer

	err = tmpl.Execute(&buf, host)
	if err != nil {
		log.Error(err)
	}

	return buf.String()
}
