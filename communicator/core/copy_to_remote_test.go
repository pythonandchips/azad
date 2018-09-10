package core

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/pythonandchips/azad/expect"
	plugintesting "github.com/pythonandchips/azad/plugin/testing"
	"github.com/stretchr/testify/assert"
	"github.com/zclconf/go-cty/cty"
)

func TestCopyToRemoteCommand(t *testing.T) {
	t.Run("with source and dest specified", func(t *testing.T) {
		t.Run("and command is successful", func(t *testing.T) {
			wd, _ := os.Getwd()
			rolepath := filepath.Join(
				wd, "fixtures",
			)
			fakeContext := plugintesting.NewFakeContext()
			fakeContext.SetRolePath(rolepath)
			fakeContext.SetVars(map[string]cty.Value{
				"source": cty.StringVal("copy_file.conf"),
				"dest":   cty.StringVal("$HOME/file.conf"),
			})
			_, err := copyToRemoteCommand(fakeContext)
			expect.NoErrors(t, err)

			command := fakeContext.CommandRan()
			expect.EqualFatal(t, len(command.Command), 4)
			assert.Equal(t, command.Command[0], `echo "dd5eb13e02ae500fb6681bf6c9300659ad06c601 $HOME/file.conf" | sha1sum -c -`)
			assert.Equal(t, command.Command[1], `if [ $? = 0 ]; then exit 40; fi`)
			assert.Equal(t, command.Command[2], "filebase64encoded=dGhpcyBpcyBhIGNvbmYgZmlsZQoKdG8gYmUgdHJhbnNmZXJlZCB0byBzZXJ2ZXIKCmFzIGlzCg==")
			assert.Equal(t, command.Command[3], "echo $filebase64encoded | base64 -d > $HOME/file.conf")
		})
	})
}
