package main

// var agentCmd = &cobra.Command{
// 	Use:     "agent",
// 	Short:   "Run a znet agent",
// 	Long:    "Subscribe to RPC events to perform actions",
// 	Example: "znet agent",
// 	// Run:     runAgent,
// }

// func init() {
// 	rootCmd.AddCommand(agentCmd)
// }

func main() {
	// 	initLogger()

	// 	z, err := znet.NewZnet(cfgFile)
	// 	if err != nil {
	// 		log.Fatalf("failed to create znet: %s", err)
	// 	}

	// 	z.Config.RPC.ServerAddress = viper.GetString("rpc.server_address")

	// 	if z.Config.RPC.ServerAddress == "" {
	// 		log.Fatal("no rpc.server configuration specified")
	// 	}

	// 	cfg := &config.Config{
	// 		Vault: z.Config.Vault,
	// 		TLS:   z.Config.TLS,
	// 	}

	// 	conn := comms.StandardRPCClient(z.Config.RPC.ServerAddress, *cfg)

	// 	if z.Config.Agent == nil {
	// 		log.Fatal("unable to create agent with nil Agent configuration")
	// 	}

	// 	var agentServer *agent.Agent
	// 	if z.Config.Agent != nil {
	// 		agentServer, err = agent.NewAgent(z.Config, conn)
	// 		if err != nil {
	// 			log.Fatal(err)
	// 		}
	// 	}

	// 	err = agentServer.Start()
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}

	// 	sigs := make(chan os.Signal, 1)
	// 	done := make(chan bool, 1)

	// 	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// 	go func() {
	// 		sig := <-sigs
	// 		log.Debugf("caught signal: %s", sig.String())

	// 		done <- true
	// 	}()

	// 	<-done

	// 	log.Debug("terminating RPC server")
	// 	err = agentServer.Stop()
	// 	if err != nil {
	// 		log.Error(err)
	// 	}

	// 	log.Debug("closing RPC client connection")
	// 	err = conn.Close()
	// 	if err != nil {
	// 		log.Error(err)
	// 	}
}
