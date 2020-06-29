package inventory

import (
	"fmt"
	"time"

	ldap "github.com/go-ldap/ldap/v3"
	log "github.com/sirupsen/logrus"
)

func (i *Inventory) SetAttribute(dn, attributeName, attributeValue string, replace bool) error {
	log.Tracef("SetAttribute: %s: %s=%s", dn, attributeName, attributeValue)

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

func (i *Inventory) UpdateTimestamp(dn string, object string) error {
	now := time.Now()

	objectName := fmt.Sprintf("%sLastSeen", object)
	return i.SetAttribute(dn, objectName, now.Format(time.RFC3339), true)
}
