package inventory

import (
	"crypto/tls"
	"fmt"
	"strconv"
	"time"

	ldap "github.com/go-ldap/ldap/v3"
	log "github.com/sirupsen/logrus"
	"github.com/xaque208/znet/internal/config"
)

// NewLDAPClient constructs an LDAP client to return.
func NewLDAPClient(cfg *config.LDAPConfig) (*ldap.Conn, error) {

	if cfg.BindDN == "" || cfg.BindPW == "" || cfg.BaseDN == "" {
		return nil, fmt.Errorf("incomplete LDAP credentials, need [BindDN, BindPW, BaseDN]")
	}

	l, err := ldap.DialTLS(
		"tcp",
		fmt.Sprintf("%s:%d", cfg.Host, 636),
		&tls.Config{InsecureSkipVerify: true},
	)
	if err != nil {
		return nil, err
	}

	l.SetTimeout(15 * time.Second)

	// First bind with a read only user
	err = l.Bind(cfg.BindDN, cfg.BindPW)
	if err != nil {
		return nil, err
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
					fmt.Sprintf("%s:%d", cfg.Host, 636),
					&tls.Config{InsecureSkipVerify: true},
				)
				if err != nil {
					log.Error(err)
				}

				err = l.Bind(cfg.BindDN, cfg.BindPW)
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
		v, err := strconv.ParseBool(string(b))
		if err != nil {
			log.Errorf("unable to parse bool: %+v", err)
		}

		values = append(values, v)
	}

	return values
}
