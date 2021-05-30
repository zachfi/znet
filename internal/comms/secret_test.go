package comms

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestVaultStuff(t *testing.T) {
	ln, client := createTestVault(t)
	defer ln.Close()

	require.NotNil(t, ln)
	require.NotNil(t, client)
}
