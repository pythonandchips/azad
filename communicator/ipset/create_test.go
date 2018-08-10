package ipset

import (
	"testing"

	"github.com/pythonandchips/azad/expect"
	plugintesting "github.com/pythonandchips/azad/plugin/testing"
)

func TestCreateCommand(t *testing.T) {
	t.Run("with only command specified", func(t *testing.T) {
		t.Run("and command is successful", func(t *testing.T) {
			fakeContext := plugintesting.NewFakeContext()
			fakeContext.SetVars(map[string]string{
				"set":   "application",
				"entry": "bitmap:port range 0-65535",
			})
			_, err := createCommand(fakeContext)
			expect.NoErrors(t, err)

			command := fakeContext.CommandRan()

			expectedCommand := []string{
				`ipset create application bitmap:port range 0-65535`,
			}

			plugintesting.AssertCommand(t, command, expectedCommand)
		})
		t.Run("with dest set", func(t *testing.T) {
			fakeContext := plugintesting.NewFakeContext()
			fakeContext.SetVars(map[string]string{
				"set":   "application",
				"entry": "bitmap:port range 0-65535",
				"dest":  "ipset-application",
			})
			_, err := createCommand(fakeContext)
			expect.NoErrors(t, err)

			command := fakeContext.CommandRan()

			expectedCommand := []string{
				`ipset create application bitmap:port range 0-65535`,
				`ipset save application > ipset-application`,
			}

			plugintesting.AssertCommand(t, command, expectedCommand)
		})
	})
}
