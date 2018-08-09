package ipset

import (
	"testing"

	"github.com/pythonandchips/azad/expect"
	plugintesting "github.com/pythonandchips/azad/plugin/testing"
)

func TestOpenPortCommand(t *testing.T) {
	t.Run("with only command specified", func(t *testing.T) {
		t.Run("and command is successful", func(t *testing.T) {
			fakeContext := plugintesting.NewFakeContext()
			fakeContext.SetVars(map[string]string{
				"set":  "application",
				"port": "6379",
			})
			_, err := openPortCommand(fakeContext)
			expect.NoErrors(t, err)

			command := fakeContext.CommandRan()

			expectedCommand := []string{
				`ipset -! add application 6379`,
			}

			plugintesting.AssertCommand(t, command, expectedCommand)
		})
		t.Run("and command is successful", func(t *testing.T) {
			fakeContext := plugintesting.NewFakeContext()
			fakeContext.SetVars(map[string]string{
				"set":  "application",
				"port": "6379",
				"dest": "ipset-application",
			})
			_, err := openPortCommand(fakeContext)
			expect.NoErrors(t, err)

			command := fakeContext.CommandRan()

			expectedCommand := []string{
				`ipset -! add application 6379`,
				`ipset save application > ipset-application`,
			}

			plugintesting.AssertCommand(t, command, expectedCommand)
		})
	})
}
