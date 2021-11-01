package netconfig

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"text/template"

	"github.com/imdario/mergo"
	"github.com/scottdware/go-junos"
	log "github.com/sirupsen/logrus"

	"github.com/xaque208/znet/modules/inventory"
)

// Host is a single configurable host.
type Host struct {
	HostName    string
	NetworkHost *inventory.NetworkHost
	Data        HostData
	Environment map[string]string
}

// NetConfig is enough data to configure some network hosts.
type NetConfig struct {
	Data      Data
	Hosts     []Host
	ConfigDir string
	junosAuth *junos.AuthMethod
}

// NewNetConfig is used to build a new *NetConfig.
func NewNetConfig(configDir string, hosts []*inventory.NetworkHost, auth *junos.AuthMethod, env map[string]string) (*NetConfig, error) {
	data, err := loadData(configDir)
	if err != nil {
		return nil, err
	}

	if len(hosts) == 0 {
		return nil, fmt.Errorf("unable to configure zero hosts")
	}

	if auth == nil {
		return nil, fmt.Errorf("unable to auth with nil")
	}

	nc := &NetConfig{
		Data:      data,
		ConfigDir: configDir,
		junosAuth: auth,
	}

	log.WithFields(log.Fields{
		"host_count": len(hosts),
	}).Debug("netconfig")

	for _, h := range hosts {
		if h == nil {
			log.Error("unable to configure nil host")
			continue
		}

		netHost := h

		host := Host{
			NetworkHost: netHost,
			HostName:    strings.Join([]string{h.Name, h.Domain}, "."),
			Environment: env,
		}

		d := nc.dataForHost(host)
		host.Data = d

		nc.Hosts = append(nc.Hosts, host)
	}

	return nc, nil
}

// LoadData receives a configuration directory from which to load the data for Znet.
func loadData(configDir string) (Data, error) {
	log.Debugf("loading data from: %s", configDir)
	dataConfig := Data{}
	err := loadYamlFile(fmt.Sprintf("%s/%s", configDir, "data.yaml"), &dataConfig)
	if err != nil {
		return Data{}, fmt.Errorf("failed to load yaml file %s: %s", configDir, err)
	}

	return dataConfig, nil
}

// ConfigureNetwork configures all discovered network devices.
func (n *NetConfig) ConfigureNetwork(commit bool, confirm int, diff bool) error {
	if n == nil {
		return fmt.Errorf("unable to configure network with nil NetConfig")
	}

	wg := sync.WaitGroup{}
	for _, host := range n.Hosts {
		wg.Add(1)
		go func(h Host) {
			if h.NetworkHost.Platform == "junos" {
				log.Debugf("configuring network host: %+v", h.HostName)

				err := n.ConfigureNetworkHost(h, commit, confirm, diff)
				if err != nil {
					log.Error(err)
				}
			}

			wg.Done()
		}(host)
	}
	wg.Wait()

	return nil
}

// ConfigureNetworkHost renders the templates using associated data for a
// network host.  The hosts about which to load the templates, are retrieved
// from LDAP.
func (n *NetConfig) ConfigureNetworkHost(host Host, commit bool, confirm int, diff bool) error {

	// log.Debugf("Using auth: %+v", auth)
	session, err := junos.NewSession(host.HostName, n.junosAuth)
	if err != nil {
		return err
	}

	defer session.Close()

	templates := n.templatesForDevice(host)
	// log.Debugf("Templates for host %s: %+v", host.Name, templates)

	var renderedTemplates []string
	for _, t := range templates {
		result := n.renderHostTemplateFile(host, t)
		renderedTemplates = append(renderedTemplates, result)
		// log.Infof("Result: %+v", result)
	}

	if diff {
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

	diffResult, err := session.Diff(0)
	if err != nil {
		return err
	}

	if len(diffResult) > 1 {
		log.Infof("configuration changes for %s: %s", host.HostName, diffResult)

		if commit {
			if confirm > 0 {
				err = session.CommitConfirm(confirm)
				if err != nil {
					return err
				}
			} else {
				err = session.Commit()
				if err != nil {
					return err
				}
			}

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

// DataForDevice returns HostData for a given NetworkHost.
func (n *NetConfig) dataForHost(host Host) HostData {
	hostData := HostData{}

	for _, f := range n.hierarchyForDevice(host) {
		fileHostData := HostData{}
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
func (n *NetConfig) hierarchyForDevice(host Host) []string {
	var files []string

	paths := templateStringsForDevice(host, n.Data.Hierarchy)

	for _, p := range paths {
		templateAbs := fmt.Sprintf("%s/%s/%s", n.ConfigDir, n.Data.DataDir, p)
		if _, err := os.Stat(templateAbs); err == nil {
			files = append(files, templateAbs)
		} else if os.IsNotExist(err) {
			log.Tracef("data file %s does not exist", templateAbs)
		}
	}

	return files
}

// templatesForDevice returns a list of template paths for a given host.
func (n *NetConfig) templatesForDevice(host Host) []string {
	var files []string

	log.Tracef("loading templates for host: %+v", host)

	paths := templateStringsForDevice(host, n.Data.TemplatePaths)

	for _, p := range paths {
		templateAbs := fmt.Sprintf("%s/%s/%s", n.ConfigDir, n.Data.TemplateDir, p)
		if _, err := os.Stat(templateAbs); err == nil {
			globPattern := fmt.Sprintf("%s/*.tmpl", templateAbs)
			foundFiles, globErr := filepath.Glob(globPattern)
			if globErr != nil {
				log.Error(globErr)
			} else {
				files = append(files, foundFiles...)
			}
		} else if os.IsNotExist(err) {
			log.Warnf("template path %s does not exist", templateAbs)
		}
	}

	log.Tracef("found %d templates for host: %s", len(files), host.HostName)

	return files
}

// RenderHostTemplateFile renders a template file using a Host object.
func (n *NetConfig) renderHostTemplateFile(host Host, path string) string {

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

	err = tmpl.Execute(&buf, host)
	if err != nil {
		log.Error(err)
	}

	return buf.String()
}
