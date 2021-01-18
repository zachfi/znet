// +build unit

package lights

import "testing"

func TestHueLight_interface(t *testing.T) {
	var l Light = hueLight{}
	t.Log(l)
}
