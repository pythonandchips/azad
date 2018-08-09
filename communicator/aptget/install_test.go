package aptget

import (
	"testing"

	"github.com/pythonandchips/azad/expect"
	plugintesting "github.com/pythonandchips/azad/plugin/testing"
	"github.com/stretchr/testify/assert"
)

func TestInstallCommand(t *testing.T) {
	t.Run("with package name", func(t *testing.T) {
		t.Run("and command is successful", func(t *testing.T) {
			fakeContext := plugintesting.NewFakeContext()
			fakeContext.SetVars(map[string]string{
				"package": "curl",
			})
			_, err := installCommand(fakeContext)
			expect.NoErrors(t, err)

			command := fakeContext.CommandRan()
			expect.EqualFatal(t, len(command.Command), 1)
			assert.Equal(t, command.Command[0], `apt-get install -yq curl`)
		})
		t.Run("and update is specified", func(t *testing.T) {
			fakeContext := plugintesting.NewFakeContext()
			fakeContext.SetVars(map[string]string{
				"package": "curl",
				"update":  "true",
			})
			_, err := installCommand(fakeContext)
			expect.NoErrors(t, err)

			command := fakeContext.CommandRan()
			expect.EqualFatal(t, len(command.Command), 2)
			assert.Equal(t, command.Command[0], `apt-get update -yq`)
			assert.Equal(t, command.Command[1], `apt-get install -yq curl`)
		})
		t.Run("and deb is specified", func(t *testing.T) {
			fakeContext := plugintesting.NewFakeContext()
			fakeContext.SetVars(map[string]string{
				"package": "erlang",
				"deb":     "https://packages.erlang-solutions.com/erlang-solutions_1.0_all.deb",
			})
			_, err := installCommand(fakeContext)
			expect.NoErrors(t, err)

			command := fakeContext.CommandRan()
			expectedCommands := []string{
				`wget -qO /tmp/erlang-solutions_1.0_all.deb https://packages.erlang-solutions.com/erlang-solutions_1.0_all.deb`,
				`dpkg -i /tmp/erlang-solutions_1.0_all.deb`,
				`rm /tmp/erlang-solutions_1.0_all.deb`,
				`apt-get update -yq`,
				`apt-get install -yq erlang`,
			}
			plugintesting.AssertCommand(t, command, expectedCommands)
		})
	})
}
