package aptget

import (
	"testing"

	"github.com/pythonandchips/azad/expect"
	plugintesting "github.com/pythonandchips/azad/plugin/testing"
	"github.com/zclconf/go-cty/cty"
)

func TestAddRepoCommand(t *testing.T) {
	t.Run("with package name", func(t *testing.T) {
		t.Run("and command is successful", func(t *testing.T) {
			fakeContext := plugintesting.NewFakeContext()
			fakeContext.SetVars(map[string]cty.Value{
				"repo": cty.StringVal("deb https://sensu.global.ssl.fastly.net/apt stretch main"),
				"key":  cty.StringVal("https://sensu.global.ssl.fastly.net/apt/pubkey.gpg"),
			})
			_, err := addRepoCommand(fakeContext)
			expect.NoErrors(t, err)

			command := fakeContext.CommandRan()
			expectedCommands := []string{
				`wget -qO - https://sensu.global.ssl.fastly.net/apt/pubkey.gpg | sudo apt-key add -`,
				`echo "deb https://sensu.global.ssl.fastly.net/apt stretch main" > /etc/apt/sources.list.d/sensu.global.ssl.fastly.net.list`,
			}
			plugintesting.AssertCommand(t, command, expectedCommands)
		})
		t.Run("with update specified", func(t *testing.T) {
			fakeContext := plugintesting.NewFakeContext()
			fakeContext.SetVars(map[string]cty.Value{
				"repo":   cty.StringVal("deb https://sensu.global.ssl.fastly.net/apt stretch main"),
				"key":    cty.StringVal("https://sensu.global.ssl.fastly.net/apt/pubkey.gpg"),
				"update": cty.BoolVal(true),
			})
			_, err := addRepoCommand(fakeContext)
			expect.NoErrors(t, err)

			command := fakeContext.CommandRan()
			expectedCommands := []string{
				`wget -qO - https://sensu.global.ssl.fastly.net/apt/pubkey.gpg | sudo apt-key add -`,
				`echo "deb https://sensu.global.ssl.fastly.net/apt stretch main" > /etc/apt/sources.list.d/sensu.global.ssl.fastly.net.list`,
				`apt-get update -qy`,
			}
			plugintesting.AssertCommand(t, command, expectedCommands)
		})
	})
}
