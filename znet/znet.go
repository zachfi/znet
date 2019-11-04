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
	ldap "gopkg.in/ldap.v2"
)

// Znet is the core object for this project.  It keeps track of the data, configuration and flow control for starting the server process.
type Znet struct {
	ConfigDir   string
	Config      Config
	Data        Data
	Environment map[string]string
	listener    *Listener
	ldapClient  *ldap.Conn
}

// NewZnet creates and returns a new Znet object.
func NewZnet(file string) (*Znet, error) {
	z := &Znet{}

	config, err := loadConfig(file)
	if err != nil {
		return &Znet{}, err
	}

	z.Config = config

	if config.LDAP.BindDN != "" && config.LDAP.BindPW != "" {
		ldapConn, err := NewLDAPClient(config.LDAP)
		if err != nil {
			return &Znet{}, fmt.Errorf("Failed LDAP connection: %s", err)
		}

		z.ldapClient = ldapConn
	}

	err = z.LoadEnvironment()
	if err != nil {
		return &Znet{}, fmt.Errorf("Failed to load environment: %s", err)
	}

	return z, nil
}

// LoadConfig receives a file path for a configuration to load.
func loadConfig(file string) (Config, error) {
	filename, _ := filepath.Abs(file)
	log.Debugf("Loading config from: %s", filename)
	config := Config{}
	err := loadYamlFile(filename, &config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}

// LoadData receives a configuration directory from which to load the data for Znet.
func (z *Znet) LoadData(configDir string) {
	log.Debugf("Loading data from: %s", configDir)
	dataConfig := Data{}
	loadYamlFile(fmt.Sprintf("%s/%s", configDir, "data.yaml"), &dataConfig)

	z.Data = dataConfig
}

// ConfigureNetworkHost renders the templates using associated data for a network host.  The hosts about which to load the templates, are retrieved from LDAP.
func (z *Znet) ConfigureNetworkHost(host *NetworkHost, commit bool, auth *junos.AuthMethod, show bool) error {

	// log.Debugf("Using auth: %+v", auth)
	session, err := junos.NewSession(host.HostName, auth)
	if err != nil {
		return err
	}

	defer session.Close()

	// log.Warnf("Auth: %+v", auth)

	// log.Warnf("Znet: %+v", z)
	// log.Warnf("Commit: %t", commit)
	// log.Warnf("Host: %+v", host)
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
		log.Debugf("RenderedTemplates: %+v", renderedTemplates)
	}

	err = session.Lock()
	if err != nil {
		return err
	}

	err = session.Config(renderedTemplates, "text", false)
	if err != nil {
		return err
	}

	diff, err := session.Diff(0)
	if err != nil {
		return err
	}

	if len(diff) > 1 {
		log.Infof("Configuration changes for %s: %s", host.HostName, diff)

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

	err = session.Unlock()
	if err != nil {
		return err
	}

	return nil
}

// TemplateStringsForDevice renders a list of template strings given a host.
func (z *Znet) TemplateStringsForDevice(host NetworkHost, templates []string) []string {
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
				files = append(files, foundFiles...)
			}

		} else if os.IsNotExist(err) {
			log.Warnf("Template path %s does not exist", templateAbs)
		}
	}

	return files
}

// RenderHostTemplateFile renders a template file using a Host object.
func (z *Znet) RenderHostTemplateFile(host NetworkHost, path string) string {
	// log.Debugf("Rendering host template file %s for host %s", path, host.Name)

	b, err := ioutil.ReadFile(path)
	if err != nil {
		log.Error(err)
	}

	str := string(b)
	tmpl, err := template.New("test").Parse(str)
	if err != nil {
		log.Error(err)
	}

	var buf bytes.Buffer

	err = tmpl.Execute(&buf, host)
	if err != nil {
		log.Error(err)
	}

	return buf.String()
}

// Close calls
func (z *Znet) Close() error {

	z.ldapClient.Close()

	return nil
}
