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

	"github.com/xaque208/znet/internal/comms"
	"github.com/xaque208/znet/internal/config"
	"github.com/xaque208/znet/internal/inventory"
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
	initLogger()

	z, err := znet.NewZnet(cfgFile)
	if err != nil {
		log.Error(err)
	}

	z.Config.RPC.ServerAddress = viper.GetString("rpc.server_address")

	if z.Config.RPC.ServerAddress == "" {
		log.Fatal("no rpc.server configuration specified")
	}

	cfg := &config.Config{
		Vault: z.Config.Vault,
		TLS:   z.Config.TLS,
	}

	conn := comms.StandardRPCClient(z.Config.RPC.ServerAddress, *cfg)

	defer func() {
		err = conn.Close()
		if err != nil {
			log.Error(err)
		}
	}()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	inventoryClient := inventory.NewInventoryClient(conn)

	stream, err := inventoryClient.ListNetworkHosts(ctx, &inventory.Empty{})
	if err != nil {
		log.Errorf("stream error: %s", err)
	}

	for {
		var d *inventory.NetworkHost

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
