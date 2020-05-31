package inventory

import (
	"crypto/tls"
	"fmt"
	"time"

	ldap "github.com/go-ldap/ldap/v3"
	log "github.com/sirupsen/logrus"
)

// NewLDAPClient constructs an LDAP client to return.
func NewLDAPClient(config LDAPConfig) (*ldap.Conn, error) {

	if config.BindDN == "" || config.BindPW == "" {
		return &ldap.Conn{}, fmt.Errorf("incomplete LDAP credentials")
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
				log.Debug("reconnecting to LDAP...")
				var err error
				l.Close()

				l, err = ldap.DialTLS(
					"tcp",
					fmt.Sprintf("%s:%d", config.Host, 636),
					&tls.Config{InsecureSkipVerify: true},
				)
				if err != nil {
					log.Error(err)
				}

				err = l.Bind(config.BindDN, config.BindPW)
				if err != nil {
					log.Error(err)
				}
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
