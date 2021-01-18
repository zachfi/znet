// +build unit

package lights

import "testing"

func TestZigbeeLight_interface(t *testing.T) {
	var l Light = zigbeeLight{}
	t.Log(l)
}
