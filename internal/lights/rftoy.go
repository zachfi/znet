package lights

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/xaque208/rftoy/rftoy"
	"github.com/xaque208/znet/internal/config"
)

type rftoyLight struct {
	config   *config.LightsConfig
	endpoint *rftoy.RFToy
}

func NewRFToyLight(cfg *config.LightsConfig) (Handler, error) {
	if cfg.RFToy == nil {
		return nil, fmt.Errorf("unable to create new rftoy light with nil config")
	}

	return rftoyLight{
		endpoint: &rftoy.RFToy{Address: cfg.RFToy.Endpoint},
	}, nil
}

func (l rftoyLight) Off(groupName string) error {
	room, err := l.config.Room(groupName)
	if err != nil {
		return err
	}

	log.WithFields(log.Fields{
		"ids": room.IDs,
	}).Debug("turning off rftoy ids")

	var errors []string
	for _, i := range room.IDs {
		err := l.endpoint.Off(i)
		if err != nil {
			errors = append(errors, err.Error())
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf(strings.Join(errors, ", "))
	}

	return nil
}

func (l rftoyLight) On(groupName string) error {
	room, err := l.config.Room(groupName)
	if err != nil {
		return err
	}

	log.WithFields(log.Fields{
		"ids": room.IDs,
	}).Debug("turning on rftoy ids")

	var errors []string
	for _, i := range room.IDs {
		err := l.endpoint.On(i)
		if err != nil {
			errors = append(errors, err.Error())
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf(strings.Join(errors, ", "))
	}

	return nil
}

func (l rftoyLight) Alert(groupName string) error {
	return nil
}

func (l rftoyLight) Toggle(groupName string) error {
	return nil
}

func (l rftoyLight) Dim(groupName string, brightness int32) error {
	return nil
}

func (l rftoyLight) SetColor(groupName string, hex string) error {
	return nil
}

func (l rftoyLight) RandomColor(groupName string, hex []string) error {
	return nil
}
