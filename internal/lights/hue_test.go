// +build unit

package lights

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHueLight_interface(t *testing.T) {
	var l Handler = &hueLight{}
	require.NotNil(t, l)
}
