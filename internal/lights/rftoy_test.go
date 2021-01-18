// +build unit

package lights

import "testing"

func TestRftoyLight_interface(t *testing.T) {
	var l Light = rftoyLight{}
	t.Log(l)
}
