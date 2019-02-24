package znet

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/alecthomas/template"
	"github.com/imdario/mergo"
	junos "github.com/scottdware/go-junos"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Znet struct {
	ConfigDir string
	Config    Config
	Data      Data
	listener  *Listener
}

func NewZnet(file string) *Znet {
	z := &Znet{}
	z.LoadConfig(file)
	return z
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
	auth := &junos.AuthMethod{
		Username:   viper.GetString("junos.username"),
		PrivateKey: viper.GetString("junos.keyfile"),
	}

	log.Debugf("Connecting to device: %s", host.HostName)
	session, err := junos.NewSession(host.HostName, auth)
	defer session.Close()
	if err != nil {
		log.Error(err)
	}

	// log.Warnf("Auth: %+v", auth)

	// log.Warnf("Znet: %+v", z)
	// log.Warnf("Commit: %t", commit)
	// log.Warnf("Host: %+v", host)
	templates := z.TemplatesForDevice(*host)
	log.Debugf("Templates for host %s: %+v", host.Name, templates)

	host.Data = z.DataForDevice(*host)
	// log.Debugf("Data: %+v", host.Data)

	var renderedTemplates []string
	for _, t := range templates {
		result := z.RenderHostTemplateFile(*host, t)
		renderedTemplates = append(renderedTemplates, result)
		// log.Infof("Result: %+v", result)
	}
	// log.Debugf("RenderedTemplates: %+v", renderedTemplates)

	err = session.Lock()
	if err != nil {
		log.Error(err)
	}

	err = session.Config(renderedTemplates, "text", false)
	if err != nil {
		log.Error(err)
	}

	diff, err := session.Diff(0)
	if err != nil {
		log.Error(err)
	}

	if len(diff) > 1 {
		log.Infof("Configuration changes for %s: %s", host.HostName, diff)

		if commit {
			err = session.Commit()
			if err != nil {
				log.Error(err)
			}
		} else {
			err = session.Config("rollback", "text", false)
			if err != nil {
				log.Error(err)
			}

		}
	}

	err = session.Unlock()
	if err != nil {
		log.Error(err)
	}

}

// TemplateStringsForDevice renders a list of template strings given a host.
func (z *Znet) TemplateStringsForDevice(host NetworkHost, templates []string) []string {
	var strings []string

	for _, t := range templates {
		tmpl, err := template.New("template").Parse(t)

		var buf bytes.Buffer

		err = tmpl.Execute(&buf, host)
		if err != nil {
			log.Error(err)
		}

		strings = append(strings, buf.String())
	}

	return strings
}

// DataForDevice
func (z *Znet) DataForDevice(host NetworkHost) HostData {
	hostData := HostData{}

	for _, f := range z.HierarchyForDevice(host) {

		fileHostData := HostData{}
		loadYamlFile(f, &fileHostData)

		if err := mergo.Merge(&hostData, fileHostData, mergo.WithOverride); err != nil {
			log.Error(err)
		}

	}

	return hostData
}

// HierarchyForDevice retuns a list of file paths to consult for the data hierarchy.
func (z *Znet) HierarchyForDevice(host NetworkHost) []string {
	var files []string

	paths := z.TemplateStringsForDevice(host, z.Data.Hierarchy)

	for _, p := range paths {
		templateAbs := fmt.Sprintf("%s/%s/%s", z.ConfigDir, z.Data.DataDir, p)
		if _, err := os.Stat(templateAbs); err == nil {
			files = append(files, templateAbs)

		} else if os.IsNotExist(err) {
			log.Warnf("Data file %s does not exist", templateAbs)
		}

	}

	return files
}

// TemplatesForDevice returns a list of template paths for a given host.
func (z *Znet) TemplatesForDevice(host NetworkHost) []string {
	var files []string

	paths := z.TemplateStringsForDevice(host, z.Data.TemplatePaths)

	for _, p := range paths {
		templateAbs := fmt.Sprintf("%s/%s/%s", z.ConfigDir, z.Data.TemplateDir, p)
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

// RenderHostTemplateFile renders a template file using a Host object.
func (z *Znet) RenderHostTemplateFile(host NetworkHost, path string) string {
	log.Debugf("Rendering host template file %s for host %s", path, host.Name)

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
