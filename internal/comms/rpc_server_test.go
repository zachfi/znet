package comms

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStandardRPCServer_interface(t *testing.T) {
	var s RPCServerFunc = StandardRPCServer
	require.NotNil(t, s)
}
