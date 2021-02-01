// +build unit

package lights

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRftoyLight_interface(t *testing.T) {
	var l Handler = rftoyLight{}
	require.NotNil(t, l)
}
