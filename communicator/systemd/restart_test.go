package systemd

import (
	"testing"

	"github.com/pythonandchips/azad/expect"
	plugintesting "github.com/pythonandchips/azad/plugin/testing"
	"github.com/stretchr/testify/assert"
	"github.com/zclconf/go-cty/cty"
)

func TestRestart(t *testing.T) {
	t.Run("restarts a systemd service", func(t *testing.T) {
		t.Run("and command is successful", func(t *testing.T) {
			fakeContext := plugintesting.NewFakeContext()
			fakeContext.SetVars(map[string]cty.Value{
				"service": cty.StringVal("networking"),
			})
			_, err := restartCommand(fakeContext)
			expect.NoErrors(t, err)

			command := fakeContext.CommandRan()
			expect.EqualFatal(t, len(command.Command), 1)
			assert.Equal(t, command.Command[0], `systemctl restart networking`)
		})
	})
}
