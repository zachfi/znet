package cmd

import (
	"context"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	junos "github.com/scottdware/go-junos"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	pb "github.com/xaque208/znet/rpc"
	"github.com/xaque208/znet/znet"
)

var arpwatchCmd = &cobra.Command{
	Use:     "arpwatch",
	Short:   "Export Junos ARP data Pometheus",
	Long:    "Run an ARP reporter",
	Run:     runArpwtach,
	Example: "znet arpwatch -v -i 10",
}

var (
	interval      int
	junosUsername string
	junosPassword string
)

func init() {
	rootCmd.AddCommand(arpwatchCmd)

	arpwatchCmd.PersistentFlags().IntVarP(&interval, "interval", "i", 60, "The interval at which to update the data")
	arpwatchCmd.PersistentFlags().StringVarP(&junosUsername, "username", "", "", "The Junos username")
	arpwatchCmd.PersistentFlags().StringVarP(&junosPassword, "password", "", "", "The Junos password")

	// err := rootCmd.MarkPersistentFlagRequired("config")
	// if err != nil {
	// 	log.Error(err)
	// }

	err := viper.BindPFlag("interval", rootCmd.PersistentFlags().Lookup("interval"))
	if err != nil {
		log.Error(err)
	}
}

func runArpwtach(cmd *cobra.Command, args []string) {
	formatter := log.TextFormatter{
		FullTimestamp: true,
	}

	log.SetFormatter(&formatter)
	if trace {
		log.SetLevel(log.TraceLevel)
	} else if verbose {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	z, err := znet.NewZnet(cfgFile)
	if err != nil {
		log.Fatal(err)
	}

	z.Config.RPC.ServerAddress = viper.GetString("rpc.server")

	conn := znet.NewConn(z.Config.RPC.ServerAddress, z.Config)
	defer func() {
		err = conn.Close()
		if err != nil {
			log.Error(err)
		}
	}()

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		log.Warnf("caught signal: %s", sig.String())

		done <- true
	}()

	telemetryClient := pb.NewTelemetryClient(conn)
	inventoryClient := pb.NewInventoryClient(conn)

	auth := &junos.AuthMethod{
		Username:   viper.GetString("junos.username"),
		PrivateKey: viper.GetString("junos.keyfile"),
	}

	scrapeJunosHost := func(wg *sync.WaitGroup, h *pb.NetworkHost) {

		hostName := strings.Join([]string{h.Name, h.Domain}, ".")
		log.Debugf("scraping ARP status from host: %s", hostName)
		session, err := junos.NewSession(hostName, auth)
		if err != nil {
			log.Error(err)
			return
		}
		defer session.Close()
		defer wg.Done()

		views, err := session.View("arp")
		if err != nil {
			log.Error(err)
			wg.Done()
			return
		}

		for _, arp := range views.Arp.Entries {
			if arp.Interface == "ppp0.0" {
				continue
			}

			log.Tracef("reporting NetworkID: %+v", arp)

			name := strings.ToLower(strings.Join([]string{arp.MACAddress, arp.Interface}, "_"))

			networkID := &pb.NetworkID{
				Name:                     name,
				IpAddress:                []string{arp.IPAddress},
				MacAddress:               []string{arp.MACAddress},
				ReportingSource:          []string{h.Name},
				ReportingSourceInterface: []string{arp.Interface},
			}

			_, err = telemetryClient.ReportNetworkID(context.Background(), networkID)
			if err != nil {
				log.Error(err)
			}

		}
	}

	log.Debugf("tick interval: %d", interval)
	ticker := time.NewTicker(time.Duration(interval) * time.Second)

	// Scrape the metrics
	go func() {
		for range ticker.C {
			wg := sync.WaitGroup{}
			resp, err := inventoryClient.ListNetworkHosts(context.Background(), &pb.Empty{})
			if err != nil {
				log.Error(err)
			}

			if resp != nil {
				for _, h := range resp.Hosts {

					if h.Platform == "" {
						continue
					}

					if h.Name == "" {
						continue
					}

					if h.Platform == "junos" {
						wg.Add(1)
						go scrapeJunosHost(&wg, h)
					}
				}
			}

			wg.Wait()
		}
	}()

	<-done
}
