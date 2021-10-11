package comms

import (
	"google.golang.org/grpc"
)

func TestRPCServer() (*grpc.Server, error) {

	return grpc.NewServer(), nil
}

// func createTestVault(t *testing.T) (net.Listener, *api.Client) {
// 	t.Helper()
//
// 	// Create an in-memory, unsealed core (the "backend", if you will).
// 	core, keyShares, rootToken := vault.TestCoreUnsealed(t)
// 	_ = keyShares
//
// 	// Start an HTTP server for the core.
// 	ln, addr := http.TestServer(t, core)
//
// 	// Create a client that talks to the server, initially authenticating with
// 	// the root token.
// 	conf := api.DefaultConfig()
// 	conf.Address = addr
//
// 	client, err := api.NewClient(conf)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	client.SetToken(rootToken)
//
// 	// Setup required secrets, policies, etc.
// 	_, err = client.Logical().Write("secret/foo", map[string]interface{}{
// 		"secret": "bar",
// 	})
// 	if err != nil {
// 		t.Fatal(err)
// 	}
//
// 	return ln, client
// }
