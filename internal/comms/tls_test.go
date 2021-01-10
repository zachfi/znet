package main

import (
	"testing"

	"github.com/johanbrandhorst/certify"
)

func TestSingletonKey(t *testing.T) {
	var x certify.KeyGenerator = &singletonKey{}
	t.Logf("singletonKey: %+v", x)
}
