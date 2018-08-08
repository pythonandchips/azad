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
	})
}
