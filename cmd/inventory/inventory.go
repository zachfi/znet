// Code generated, do not edit
package cmd

import (
// "context"
// "io"

// prompt "github.com/c-bata/go-prompt"
// log "github.com/sirupsen/logrus"
/* "github.com/spf13/cobra" */
// "github.com/spf13/viper"
// "google.golang.org/grpc/codes"
// "google.golang.org/grpc/status"

// "github.com/xaque208/znet/znet"
// "github.com/xaque208/znet/internal/comms"
// "github.com/xaque208/znet/internal/config"
// "github.com/xaque208/znet/modules/inventory"
)

/* var inventoryCommand = &cobra.Command{ */
/* 	Use:     "inventory", */
/* 	Short:   "Interact with inventory data", */
/* 	Long:    "Interact with inventory data", */
/* 	Example: "znet inv", */
/* 	// Run:     runInteractive, */
/* } */

/* var listCmd = &cobra.Command{ */
/* 	Use:     "list", */
/* 	Short:   "List items from the inventory", */
/* 	Long:    "List items from the inventory", */
/* 	Example: "znet inventory list", */
/* 	// Run:     runList, */
/* } */

// func runInteractive(cmd *cobra.Command, args []string) {
// 	initLogger()
//
// 	cfg, err := config.LoadConfig(cfgFile)
// 	if err != nil {
// 		log.Error(err)
// 	}
//
// 	inv, err := inventory.NewLDAPInventory(cfg.LDAP)
// 	if err != nil {
// 		log.Error(err)
// 	}
//
// 	completer := &inventory.InventoryInteractive{
// 		Inventory: inv,
// 	}
//
// 	p := prompt.New(
// 		completer.Executor,
// 		completer.Completer,
// 		prompt.OptionTitle("znet inv"),
// 	)
//
// 	p.Run()
// }
//
//
//
//
//
//
//
//
//
//
//
// // typeName: network_host
// // commandName: NetworkHost
//
// var listNetworkHostCmd = &cobra.Command{
// 	Use:     "network_host",
// 	Short:   "Manage network_host inventory resources",
// 	// Long:    "Run an inventory report",
// 	Example: "znet inventory network_host",
// 	Run:     runListNetworkHost,
// }
//
// func runListNetworkHost(cmd *cobra.Command, args []string) {
// 	initLogger()
//
// 	z, err := znet.NewZnet(cfgFile)
// 	if err != nil {
// 		log.Error(err)
// 	}
//
// 	z.Config.RPC.ServerAddress = viper.GetString("rpc.server_address")
//
// 	if z.Config.RPC.ServerAddress == "" {
// 		log.Fatal("no rpc.server configuration specified")
// 	}
//
// 	cfg := &config.Config{
// 		Vault: z.Config.Vault,
// 		TLS:   z.Config.TLS,
// 	}
//
// 	conn := comms.StandardRPCClient(z.Config.RPC.ServerAddress, *cfg)
//
// 	defer func() {
// 		err = conn.Close()
// 		if err != nil {
// 			log.Error(err)
// 		}
// 	}()
//
// 	ctx, cancel := context.WithCancel(context.Background())
// 	defer cancel()
//
// 	inventoryClient := inventory.NewInventoryClient(conn)
//
// 	stream, err := inventoryClient.ListNetworkHosts(ctx, &inventory.Empty{})
// 	if err != nil {
//     log.Errorf("stream error: %s", err)
// 	}
//
// 	for {
// 		var d *inventory.NetworkHost
//
// 		d, err = stream.Recv()
// 		if err != nil {
// 			if err == io.EOF {
// 				break
// 			}
//
// 			switch status.Code(err) {
// 			case codes.OK:
// 				continue
// 			default:
// 				log.Errorf("default status.Code: %+v", status.Code(err))
//         break
// 			}
// 		}
//
// 		if d != nil {
// 			log.Debugf("NetworkHost: %+v", d)
// 		}
//   }
//
// }
//
//
//
//
//
//
//
//
//
//
//
//
//
// // typeName: zigbee_device
// // commandName: ZigbeeDevice
//
// var listZigbeeDeviceCmd = &cobra.Command{
// 	Use:     "zigbee_device",
// 	Short:   "Manage zigbee_device inventory resources",
// 	// Long:    "Run an inventory report",
// 	Example: "znet inventory zigbee_device",
// 	Run:     runListZigbeeDevice,
// }
//
// func runListZigbeeDevice(cmd *cobra.Command, args []string) {
// 	initLogger()
//
// 	z, err := znet.NewZnet(cfgFile)
// 	if err != nil {
// 		log.Error(err)
// 	}
//
// 	z.Config.RPC.ServerAddress = viper.GetString("rpc.server_address")
//
// 	if z.Config.RPC.ServerAddress == "" {
// 		log.Fatal("no rpc.server configuration specified")
// 	}
//
// 	cfg := &config.Config{
// 		Vault: z.Config.Vault,
// 		TLS:   z.Config.TLS,
// 	}
//
// 	conn := comms.StandardRPCClient(z.Config.RPC.ServerAddress, *cfg)
//
// 	defer func() {
// 		err = conn.Close()
// 		if err != nil {
// 			log.Error(err)
// 		}
// 	}()
//
// 	ctx, cancel := context.WithCancel(context.Background())
// 	defer cancel()
//
// 	inventoryClient := inventory.NewInventoryClient(conn)
//
// 	stream, err := inventoryClient.ListZigbeeDevices(ctx, &inventory.Empty{})
// 	if err != nil {
//     log.Errorf("stream error: %s", err)
// 	}
//
// 	for {
// 		var d *inventory.ZigbeeDevice
//
// 		d, err = stream.Recv()
// 		if err != nil {
// 			if err == io.EOF {
// 				break
// 			}
//
// 			switch status.Code(err) {
// 			case codes.OK:
// 				continue
// 			default:
// 				log.Errorf("default status.Code: %+v", status.Code(err))
//         break
// 			}
// 		}
//
// 		if d != nil {
// 			log.Debugf("NetworkHost: %+v", d)
// 		}
//   }
//
// }
//
//
//
//
//
// func init() {
// 	inventoryCommand.AddCommand(listCmd)
//
//
//
//
//
//
//
//
// 	listCmd.AddCommand(listNetworkHostCmd)
//
//
//
//
//
//
//
//
//
//
//
//
// 	listCmd.AddCommand(listZigbeeDeviceCmd)
//
//
//
//
//
// 	// invCmd.PersistentFlags().StringVarP(&rpcServer, "rpc", "r", ":8800", "Specify RPC server address")
// 	// invCmd.Flags().BoolP("verbose", "v", false, "Raise verbosity")
//
// 	// invCmd.PersistentFlags().StringVarP(&adopt, "adopt", "a", "", "Adopt an unknown host by MAC address")
// }
