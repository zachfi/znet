package znet

import (
	"crypto/tls"
	"fmt"
	"time"

	ldap "github.com/go-ldap/ldap"
	log "github.com/sirupsen/logrus"
)

// NewLDAPClient constructs an LDAP client to return.
func NewLDAPClient(config LDAPConfig) (*ldap.Conn, error) {

	if config.BindDN == "" || config.BindPW == "" {
		return &ldap.Conn{}, fmt.Errorf("Incomplete LDAP credentials")
	}

	l, err := ldap.DialTLS(
		"tcp",
		fmt.Sprintf("%s:%d", config.Host, 636),
		&tls.Config{InsecureSkipVerify: true},
	)
	if err != nil {
		return &ldap.Conn{}, err
	}
	// defer l.Close()

	// First bind with a read only user
	err = l.Bind(config.BindDN, config.BindPW)
	if err != nil {
		return &ldap.Conn{}, err
	}

	// Handle reconnection
	go func() {
		t := time.NewTicker(2 * time.Second)
		for {
			<-t.C

			if l.IsClosing() {
				log.Debugf("Old L is closing: %+v", l)
				log.Debug("LDAP Reconnecting...")
				newl, err := ldap.DialTLS(
					"tcp",
					fmt.Sprintf("%s:%d", config.Host, 636),
					&tls.Config{InsecureSkipVerify: true},
				)
				if err != nil {
					log.Error(err)
				}

				err = newl.Bind(config.BindDN, config.BindPW)
				if err != nil {
					log.Error(err)
				}

				log.Debugf("New L is: %+v", newl)

				*l = *newl
			}
		}
	}()

	return l, nil
}

func stringValues(a *ldap.EntryAttribute) []string {
	var values []string

	for _, b := range a.ByteValues {
		values = append(values, string(b))
	}

	return values
}

func boolValues(a *ldap.EntryAttribute) []bool {
	var values []bool

	for _, b := range a.ByteValues {
		fmt.Printf("Received bool: %b", b)
		// values = append(values, b.(bool))
	}

	return values
}
