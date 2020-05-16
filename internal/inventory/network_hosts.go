package inventory

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strings"

	ldap "github.com/go-ldap/ldap"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"

	"github.com/xaque208/znet/pkg/netconfig"
)

// NetworkHost is a device that connects to the network.
type NetworkHost struct {
	Data        netconfig.HostData
	Description string
	DeviceType  string
	DN          string
	Domain      string
	Environment map[string]string
	Group       string
	HostName    string
	MACAddress  []string
	Name        string
	Platform    string
	Role        string
	Watch       bool
}

var defaultHostAttributes = []string{
	"cn",
	"dn",
	"macAddress",
	"netHostDescription",
	"netHostDomain",
	"netHostGroup",
	"netHostName",
	"netHostPlatform",
	"netHostRole",
	"netHostType",
	"netHostWatch",
}

// RecordUnknownHost stores an IP and MAC with a name to LDAP.
func (i *Inventory) RecordUnknownHost(baseDN string, address string, mac string) error {

	cn := strings.Replace(mac, ":", "", -1)

	searchRequest := ldap.NewSearchRequest(
		baseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(&(objectClass=unknownNetHost)(cn=%s))", cn),
		[]string{"cn"},
		nil,
	)

	log.Debugf("Searching LDAP with query: %s", searchRequest.Filter)

	sr, err := i.ldapClient.Search(searchRequest)
	if err != nil {
		return err
	}

	if len(sr.Entries) > 0 {
		log.Debugf("Host mac %s is already unknown", mac)
		return nil
	}

	log.Debugf("Recording unknown host %s", mac)

	dn := fmt.Sprintf("cn=%s,%s", cn, baseDN)

	a := ldap.NewAddRequest(dn, nil)
	a.Attribute("objectClass", []string{"unknownNetHost", "top"})
	a.Attribute("cn", []string{cn})
	a.Attribute("v4Address", []string{address})
	a.Attribute("macAddress", []string{mac})
	err = i.ldapClient.Add(a)
	if err != nil {
		log.Errorf("%+v", a)
		return err
	}

	return nil
}

// Update should upgrade a network host.
func (h *NetworkHost) Update() (*ssh.Conn, error) {

	sshConfig := &ssh.ClientConfig{
		User: "zach",
		Auth: []ssh.AuthMethod{
			PublicKeyFile("/home/zach/.ssh/id_ed25519"),
			SSHAgent(),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	log.Warnf("%+v", sshConfig)

	connection, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", h.HostName, 22), sshConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to dial: %s", err)
	}

	// session, err := connection.NewSession()
	// if err != nil {
	// 	return nil, fmt.Errorf("Failed to create session: %s", err)
	// }
	//
	// stdin, err := session.StdinPipe()
	// if err != nil {
	// 	return nil, fmt.Errorf("Unable to setup stdin for session: %v", err)
	// }
	// go io.Copy(stdin, os.Stdin)
	//
	// stdout, err := session.StdoutPipe()
	// if err != nil {
	// 	return nil, fmt.Errorf("Unable to setup stdout for session: %v", err)
	// }
	// go io.Copy(os.Stdout, stdout)
	//
	// stderr, err := session.StderrPipe()
	// if err != nil {
	// 	return nil, fmt.Errorf("Unable to setup stderr for session: %v", err)
	// }
	// go io.Copy(os.Stderr, stderr)
	//
	// err = session.Run("ls -l")
	// if err != nil {
	// 	log.Error(err)
	// }
	//
	// out, err := session.CombinedOutput("ls -l")
	// if err != nil {
	// 	log.Error(err)
	// }
	//
	// log.Info(out)
	log.Debugf("Connection: %+v", connection)

	return nil, nil
}

// SSHAgent builds the AuthMethod for SSH.
func SSHAgent() ssh.AuthMethod {
	if sshAgent, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK")); err == nil {
		return ssh.PublicKeysCallback(agent.NewClient(sshAgent).Signers)
	}
	return nil
}

// PublicKeyFile builds the AuthMethod for SSH.
func PublicKeyFile(file string) ssh.AuthMethod {
	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		return nil
	}

	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		return nil
	}
	return ssh.PublicKeys(key)
}
