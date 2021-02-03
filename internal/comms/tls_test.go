// +build unit

package comms

import (
	"testing"

	"github.com/johanbrandhorst/certify"
	"github.com/stretchr/testify/require"
)

func TestSingletonKey(t *testing.T) {
	var x certify.KeyGenerator = &singletonKey{}
	require.NotNil(t, x)
}
