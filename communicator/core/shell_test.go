package core

import (
	"testing"

	"github.com/pythonandchips/azad/expect"
	plugintesting "github.com/pythonandchips/azad/plugin/testing"
	"github.com/stretchr/testify/assert"
	"github.com/zclconf/go-cty/cty"
)

func TestShell(t *testing.T) {
	t.Run("with only command specified", func(t *testing.T) {
		t.Run("and command is successful", func(t *testing.T) {
			fakeContext := plugintesting.NewFakeContext()
			fakeContext.SetVars(map[string]cty.Value{
				"command": cty.StringVal("ls /home"),
			})
			_, err := shell(fakeContext)
			expect.NoErrors(t, err)

			command := fakeContext.CommandRan()
			expect.EqualFatal(t, len(command.Command), 1)
			assert.Equal(t, command.Command[0], "ls /home")
			assert.Equal(t, command.Interpreter, "sh")
		})
	})
}
