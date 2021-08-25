//go:build unit

package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRootCommand(t *testing.T) {
	assert.NotEmptyf(t, rootCmd.Use, "Need to set rootCmd.%s on rootCmd %s", "Use", rootCmd.CalledAs())
	assert.NotEmptyf(t, rootCmd.Short, "Need to set rootCmd.%s on rootCmd %s", "Short", rootCmd.CalledAs())

	for _, c := range rootCmd.Commands() {
		assert.NotEmptyf(t, c.Use, "Need to set Command.%s on Command: %s", "Use", c.CommandPath())
		assert.NotEmptyf(t, c.Short, "Need to set Command.%s on Command: %s", "Short", c.CommandPath())
		assert.NotEmptyf(t, c.Long, "Need to set Command.%s on Command: %s", "Long", c.CommandPath())
		assert.NotEmptyf(t, c.Example, "Need to set Command.%s on Command: %s", "Example", c.CommandPath())
	}
}
