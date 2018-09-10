package core

import (
	"testing"

	"github.com/pythonandchips/azad/expect"
	plugintesting "github.com/pythonandchips/azad/plugin/testing"
	"github.com/zclconf/go-cty/cty"
)

func TestRemoveFile(t *testing.T) {
	t.Run("with only command specified", func(t *testing.T) {
		t.Run("and command is successful", func(t *testing.T) {
			fakeContext := plugintesting.NewFakeContext()
			fakeContext.SetVars(map[string]cty.Value{
				"path": cty.StringVal("/path/to/file"),
			})
			_, err := removeFileCommand(fakeContext)
			expect.NoErrors(t, err)

			command := fakeContext.CommandRan()

			expectedCommand := []string{
				"rm /path/to/file",
			}
			plugintesting.AssertCommand(t, command, expectedCommand)
		})
	})
}
