package ipset

import (
	"testing"

	"github.com/pythonandchips/azad/expect"
	plugintesting "github.com/pythonandchips/azad/plugin/testing"
)

func TestSaveCommand(t *testing.T) {
	t.Run("with only command specified", func(t *testing.T) {
		t.Run("and command is successful", func(t *testing.T) {
			fakeContext := plugintesting.NewFakeContext()
			fakeContext.SetVars(map[string]string{
				"set":  "application",
				"dest": "ipset-application",
			})
			_, err := saveCommand(fakeContext)
			expect.NoErrors(t, err)

			command := fakeContext.CommandRan()

			expectedCommand := []string{
				`ipset save application > ipset-application`,
			}

			plugintesting.AssertCommand(t, command, expectedCommand)
		})
	})
}
