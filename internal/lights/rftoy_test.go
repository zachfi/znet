// +build unit

package lights

import "testing"

func TestRftoyLight_interface(t *testing.T) {
	var l Handler = rftoyLight{}
	t.Log(l)
}
