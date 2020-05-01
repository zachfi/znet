package znet

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/alecthomas/template"
	"github.com/imdario/mergo"
	"github.com/prometheus/client_golang/prometheus"
	junos "github.com/scottdware/go-junos"
	log "github.com/sirupsen/logrus"

	"github.com/xaque208/znet/internal/agent"
	"github.com/xaque208/znet/internal/inventory"
	"github.com/xaque208/znet/internal/lights"
	"github.com/xaque208/znet/pkg/events"
	"github.com/xaque208/znet/pkg/netconfig"
)

// Znet is the core object for this project.  It keeps track of the data,
// configuration and flow control for starting the server process.
type Znet struct {
	ConfigDir   string
	Config      Config
	Data        netconfig.Data
	Environment map[string]string
	Inventory   *inventory.Inventory
	Lights      *lights.Lights
}

// NewZnet creates and returns a new Znet object.
func NewZnet(file string) (*Znet, error) {
	config, err := loadConfig(file)
	if err != nil {
		return &Znet{}, fmt.Errorf("failed to load config file %s: %s", file, err)
	}

	e, err := GetEnvironmentConfig(config.Environments, "common")
	if err != nil {
		log.Error(err)
	}

	environment, err := LoadEnvironment(config.Vault, e)
	if err != nil {
		log.Errorf("Failed to load environment: %s", err)
	}

	inv := inventory.NewInventory(config.LDAP)

	lights := lights.NewLights(config.Lights)

	z := &Znet{
		Config:      config,
		Environment: environment,
		Inventory:   inv,
		Lights:      lights,
	}

	return z, nil
}

// LoadConfig receives a file path for a configuration to load.
func loadConfig(file string) (Config, error) {
	filename, _ := filepath.Abs(file)
	log.Debugf("loading config from: %s", filename)
	config := Config{}
	err := loadYamlFile(filename, &config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}

// LoadData receives a configuration directory from which to load the data for Znet.
func (z *Znet) LoadData(configDir string) {
	log.Debugf("loading data from: %s", configDir)
	dataConfig := netconfig.Data{}
	err := loadYamlFile(fmt.Sprintf("%s/%s", configDir, "data.yaml"), &dataConfig)
	if err != nil {
		log.Errorf("failed to load yaml file %s: %s", configDir, err)
	}

	z.Data = dataConfig
}

// ConfigureNetworkHost renders the templates using associated data for a network host.  The hosts about which to load the templates, are retrieved from LDAP.
func (z *Znet) ConfigureNetworkHost(host *inventory.NetworkHost, commit bool, auth *junos.AuthMethod, show bool) error {

	// log.Debugf("Using auth: %+v", auth)
	session, err := junos.NewSession(host.HostName, auth)
	if err != nil {
		return err
	}

	defer session.Close()

	templates := z.TemplatesForDevice(*host)
	// log.Debugf("Templates for host %s: %+v", host.Name, templates)

	host.Data = z.DataForDevice(*host)
	// log.Debugf("Data: %+v", host.Data)

	var renderedTemplates []string
	for _, t := range templates {
		result := z.RenderHostTemplateFile(*host, t)
		renderedTemplates = append(renderedTemplates, result)
		// log.Infof("Result: %+v", result)
	}

	if show {
		log.Debugf("renderedTemplates: %+v", renderedTemplates)
	}

	err = session.Lock()
	if err != nil {
		return fmt.Errorf("unable to lock session on %s: %s", host.HostName, err)
	}

	defer func() {
		err = session.Unlock()
		if err != nil {
			log.Errorf("error unlocking session on %s: %s", host.HostName, err)
		}
	}()

	err = session.Config(renderedTemplates, "text", false)
	if err != nil {
		return fmt.Errorf("unable to load configuration on %s: %s", host.HostName, err)
	}

	diff, err := session.Diff(0)
	if err != nil {
		return err
	}

	if len(diff) > 1 {
		log.Infof("configuration changes for %s: %s", host.HostName, diff)

		if commit {
			err = session.Commit()
			if err != nil {
				return err
			}
		} else {
			err = session.Config("rollback", "text", false)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// TemplateStringsForDevice renders a list of template strings given a host.
func (z *Znet) TemplateStringsForDevice(host inventory.NetworkHost, templates []string) []string {
	var strings []string

	for _, t := range templates {
		tmpl, err := template.New("template").Parse(t)
		if err != nil {
			log.Error(err)
		}

		var buf bytes.Buffer

		err = tmpl.Execute(&buf, host)
		if err != nil {
			log.Error(err)
		}

		strings = append(strings, buf.String())
	}

	return strings
}

// DataForDevice returns HostData for a given NetworkHost.
func (z *Znet) DataForDevice(host inventory.NetworkHost) netconfig.HostData {
	hostData := netconfig.HostData{}

	for _, f := range z.HierarchyForDevice(host) {

		fileHostData := netconfig.HostData{}
		err := loadYamlFile(f, &fileHostData)
		if err != nil {
			log.Error(err)
		}

		if err := mergo.Merge(&hostData, fileHostData, mergo.WithOverride); err != nil {
			log.Error(err)
		}
	}

	return hostData
}

// HierarchyForDevice returns a list of file paths to consult for the data hierarchy.
func (z *Znet) HierarchyForDevice(host inventory.NetworkHost) []string {
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
func (z *Znet) TemplatesForDevice(host inventory.NetworkHost) []string {
	var files []string

	paths := z.TemplateStringsForDevice(host, z.Data.TemplatePaths)

	for _, p := range paths {
		templateAbs := fmt.Sprintf("%s/%s/%s", z.ConfigDir, z.Data.TemplateDir, p)
		if _, err := os.Stat(templateAbs); err == nil {
			globPattern := fmt.Sprintf("%s/*.tmpl", templateAbs)
			foundFiles, globErr := filepath.Glob(globPattern)
			if globErr != nil {
				log.Error(globErr)
			} else {
				files = append(files, foundFiles...)
			}

		} else if os.IsNotExist(err) {
			log.Warnf("Template path %s does not exist", templateAbs)
		}
	}

	return files
}

// RenderHostTemplateFile renders a template file using a Host object.
func (z *Znet) RenderHostTemplateFile(host inventory.NetworkHost, path string) string {
	// log.Debugf("Rendering host template file %s for host %s", path, host.Name)

	b, err := ioutil.ReadFile(path)
	if err != nil {
		log.Errorf("failed read path: %s", err)
	}

	str := string(b)
	tmpl, err := template.New("test").Parse(str)
	if err != nil {
		log.Errorf("failed to parse template %s: %s", path, err)
	}

	var buf bytes.Buffer

	// Attach the znet Environment to the host
	host.Environment = z.Environment

	err = tmpl.Execute(&buf, host)
	if err != nil {
		log.Error(err)
	}

	return buf.String()
}

// Shutdown the znet connections
func (z *Znet) Stop() error {
	var err error
	// z.Inventory.ldapClient.Close()

	return err
}

// Subscriptions is yet to be used, but conforms to the interface for
// generating consumers of named events.
func (z *Znet) Subscriptions() map[string][]events.Handler {
	s := events.NewSubscriptions()

	s.Subscribe("ExecutionResult", z.executionResultHandler)

	return s.Table
}

func (z *Znet) executionResultHandler(name string, payload events.Payload) error {

	var x agent.ExecutionResult

	err := json.Unmarshal(payload, &x)
	if err != nil {
		return fmt.Errorf("failed to unmarshal %T: %s", x, err)
	}

	executionExitStatus.With(prometheus.Labels{
		"command": x.Command,
	}).Set(float64(x.ExitCode))

	executionDuration.With(prometheus.Labels{
		"command": x.Command,
	}).Set(float64(x.Duration))

	return nil
}
