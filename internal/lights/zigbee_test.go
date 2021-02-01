// +build unit

package lights

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestZigbeeLight_interface(t *testing.T) {
	var l Handler = zigbeeLight{}
	require.NotNil(t, l)
}
