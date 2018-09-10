package arpwatch

import (
	"fmt"
	"time"

	"github.com/prometheus/common/log"
	junos "github.com/scottdware/go-junos"
)

const (
	macsTable = "macs"
)

type ArpWatch struct {
	Hosts []string
	Auth  *junos.AuthMethod
}

func (a ArpWatch) Update() {
	redisClient := NewRedisClient()
	defer redisClient.Close()

	for {
		select {
		default:

			for _, h := range a.Hosts {
				session, err := junos.NewSession(h, a.Auth)
				if err != nil {
					log.Fatal(err)
				}

				views, err := session.View("arp")
				if err != nil {
					log.Fatal(err)
				}

				for _, arp := range views.Arp.Entries {
					result, err := redisClient.SIsMember(macsTable, arp.MACAddress).Result()
					if err != nil {
						log.Error(err)
					}

					if result == false {
						log.Infof("New MACAddress seen: %+v", arp.MACAddress)
						_, err := redisClient.SAdd(macsTable, arp.MACAddress).Result()
						if err != nil {
							log.Error(err)
						}
					}

					keyName := fmt.Sprintf("mac:%s", arp.MACAddress)
					redisClient.HSet(keyName, "mac", arp.MACAddress)
					redisClient.HSet(keyName, "ip", arp.IPAddress)
					redisClient.Expire(keyName, 900*time.Second)

				}

			}

			log.Debugf("Sleeping")
			time.Sleep(time.Second * 10)
		}
	}
}
