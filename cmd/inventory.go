// Code generated, do not edit
package cmd

import (
	"context"
	"io"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/xaque208/znet/rpc"
	"github.com/xaque208/znet/znet"
)

var inventoryCommand = &cobra.Command{
	Use:     "inventory",
	Short:   "Report on inventory",
	Long:    "Run an inventory report",
	Example: "znet inv",
	// Run:     runInv,
}

// typeName: network_host
// commandName: NetworkHost

var networkhostCmd = &cobra.Command{
	Use:   "network_host",
	Short: "Manage network_host inventory resources",
	// Long:    "Run an inventory report",
	Example: "znet inventory network_host",
	Run:     runNetworkHost,
}

func runNetworkHost(cmd *cobra.Command, args []string) {
	formatter := log.TextFormatter{
		DisableQuote:     true,
		DisableTimestamp: true,
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
		log.Error(err)
	}

	z.Config.RPC.ServerAddress = viper.GetString("rpc.server_address")

	if z.Config.RPC.ServerAddress == "" {
		log.Fatal("no rpc.server configuration specified")
	}

	conn := znet.NewConn(z.Config.RPC.ServerAddress, z.Config)

	defer func() {
		err = conn.Close()
		if err != nil {
			log.Error(err)
		}

		z.Stop()
	}()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	inventoryClient := rpc.NewInventoryClient(conn)

	stream, err := inventoryClient.ListNetworkHosts(ctx, &rpc.Empty{})
	if err != nil {
		log.Errorf("stream error: %s", err)
	}

	for {
		var d *rpc.NetworkHost

		d, err = stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}

			switch status.Code(err) {
			case codes.OK:
				continue
			default:
				log.Errorf("default status.Code: %+v", status.Code(err))
				break
			}
		}

		if d != nil {
			log.Debugf("NetworkHost: %+v", d)
		}
	}

}

func init() {
	inventoryCommand.AddCommand(networkhostCmd)

	// invCmd.PersistentFlags().StringVarP(&rpcServer, "rpc", "r", ":8800", "Specify RPC server address")
	// invCmd.Flags().BoolP("verbose", "v", false, "Raise verbosity")

	// invCmd.PersistentFlags().StringVarP(&adopt, "adopt", "a", "", "Adopt an unknown host by MAC address")
}