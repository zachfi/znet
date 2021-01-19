// +build unit

package lights

import "testing"

func TestZigbeeLight_interface(t *testing.T) {
	var l Handler = zigbeeLight{}
	t.Log(l)
}
