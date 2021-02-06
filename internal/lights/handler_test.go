package lights

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMockHandler_interface(t *testing.T) {
	var l Handler = &MockLight{}
	require.NotNil(t, l)
}
