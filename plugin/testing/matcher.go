package testing

import (
	"testing"

	"github.com/pythonandchips/azad/expect"
	"github.com/pythonandchips/azad/plugin"
	"github.com/stretchr/testify/assert"
)

// AssertCommand is a helper to test commands produce the expected
// shell script
func AssertCommand(t *testing.T, command plugin.Command, expected []string) {
	expect.EqualFatal(t, len(command.Command), len(expected))

	for i, command := range command.Command {
		assert.Equal(t, command, expected[i])
	}
}
