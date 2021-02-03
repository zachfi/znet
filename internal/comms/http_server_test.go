// +build unit

package comms

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStandardHTTPServer_interface(t *testing.T) {
	var s HTTPServerFunc = StandardHTTPServer
	require.NotNil(t, s)
}
