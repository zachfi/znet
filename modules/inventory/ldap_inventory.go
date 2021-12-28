package inventory

import (
	"context"
	"crypto/tls"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	ldap "github.com/go-ldap/ldap/v3"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

// LDAPInventory holds the coniguration and clients
// necessary to retrieve information from data sources.
type LDAPInventory struct {
	cfg *Config

	logger log.Logger
	tracer trace.Tracer

	ldapClient *ldap.Conn
	mux        *sync.Mutex
}

// NewLDAPInventory returns a new Inventory object from the
// received config.
func NewLDAPInventory(cfg Config, logger log.Logger) (*LDAPInventory, error) {
	ldapClient, err := NewLDAPClient(cfg, logger)
	if err != nil {
		return nil, errors.Wrap(err, "failed LDAP connection")
	}

	i := &LDAPInventory{
		cfg:        &cfg,
		logger:     logger,
		ldapClient: ldapClient,
		tracer:     otel.Tracer("harvester"),
	}

	return i, nil
}

func (i *LDAPInventory) Stop() {
	i.ldapClient.Close()
}

func (i *LDAPInventory) reconnect() error {
	// Make sure old connection if definitely closed
	i.mux.Lock()
	defer i.mux.Unlock()
	i.ldapClient.Close()

	// Connect to LDAP
	l, err := ldap.DialTLS(
		"tcp",
		fmt.Sprintf("%s:%d", i.cfg.LDAP.Host, 636),
		&tls.Config{InsecureSkipVerify: true},
	)
	if err != nil {
		return err
	}

	l.SetTimeout(15 * time.Second)

	// First bind with a read only user
	err = l.Bind(i.cfg.LDAP.BindDN, i.cfg.LDAP.BindPW)
	if err != nil {
		return err
	}

	i.ldapClient = l
	return nil
}

func (i *LDAPInventory) SetAttribute(ctx context.Context, dn, attributeName, attributeValue string, replace bool) error {
	_, span := i.tracer.Start(ctx, "SetAttribute")
	defer span.End()

	modify := ldap.NewModifyRequest(dn, nil)
	if replace {
		modify.Replace(attributeName, []string{attributeValue})
	} else {
		modify.Add(attributeName, []string{attributeValue})
	}

	err := i.ldapClient.Modify(modify)
	if err != nil {
		return err
	}

	return nil
}

func (i *LDAPInventory) UpdateTimestamp(ctx context.Context, dn string, object string) error {

	_ = level.Debug(i.logger).Log("msg", "before", "traceID", trace.SpanContextFromContext(ctx).TraceID().String())

	spanCtx, span := i.tracer.Start(ctx, "UpdateTimestamp")
	defer span.End()

	_ = level.Debug(i.logger).Log("msg", "after", "traceID", trace.SpanContextFromContext(spanCtx).TraceID().String())

	now := time.Now()

	objectName := fmt.Sprintf("%sLastSeen", object)
	return i.SetAttribute(ctx, dn, objectName, now.Format(time.RFC3339), true)
}

// NewLDAPClient constructs an LDAP client to return.
func NewLDAPClient(cfg Config, logger log.Logger) (*ldap.Conn, error) {
	logger = log.With(logger, "ldap", "client")

	if cfg.LDAP.BindDN == "" || cfg.LDAP.BindPW == "" || cfg.LDAP.BaseDN == "" {
		return nil, fmt.Errorf("incomplete LDAP credentials, need [BindDN, BindPW, BaseDN]")
	}

	l, err := ldap.DialTLS(
		"tcp",
		fmt.Sprintf("%s:%d", cfg.LDAP.Host, 636),
		&tls.Config{InsecureSkipVerify: true},
	)
	if err != nil {
		return nil, err
	}

	l.SetTimeout(15 * time.Second)

	// First bind with a read only user
	err = l.Bind(cfg.LDAP.BindDN, cfg.LDAP.BindPW)
	if err != nil {
		return nil, err
	}

	// Handle reconnection
	go func() {
		t := time.NewTicker(2 * time.Second)
		for {
			<-t.C

			if l.IsClosing() {
				_ = level.Debug(logger).Log("msg", "reconnecting to LDAP")
				var err error
				l.Close()

				l, err = ldap.DialTLS(
					"tcp",
					fmt.Sprintf("%s:%d", cfg.LDAP.Host, 636),
					&tls.Config{InsecureSkipVerify: true},
				)
				if err != nil {
					_ = level.Error(logger).Log("err", err)
				}

				err = l.Bind(cfg.LDAP.BindDN, cfg.LDAP.BindPW)
				if err != nil {
					_ = level.Error(logger).Log("err", err)
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

func boolValues(a *ldap.EntryAttribute, logger log.Logger) []bool {
	var values []bool

	for _, b := range a.ByteValues {
		v, err := strconv.ParseBool(string(b))
		if err != nil {
			_ = level.Error(logger).Log("err", err)
		}

		values = append(values, v)
	}

	return values
}
