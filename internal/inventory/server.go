package inventory

import "github.com/xaque208/znet/internal/config"

// NewInventoryServer is used to return a new InventoryServer, which implements the inventory RPC server.
func NewInventoryServer(cfg *config.LDAPConfig) (*InventoryServer, error) {
	inv, err := NewInventory(cfg)
	if err != nil {
		return nil, err
	}

	return &InventoryServer{
		inventory: inv,
	}, nil
}
