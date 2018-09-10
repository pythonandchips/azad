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

func TestTemplateCommand(t *testing.T) {
	t.Run("with source and dest specified", func(t *testing.T) {
		t.Run("and command is successful", func(t *testing.T) {
			wd, _ := os.Getwd()
			rolepath := filepath.Join(
				wd, "fixtures",
			)
			fakeContext := plugintesting.NewFakeContext()
			fakeContext.SetRolePath(rolepath)
			fakeContext.SetVars(map[string]cty.Value{
				"source": cty.StringVal("template_file.conf"),
				"dest":   cty.StringVal("$HOME/file.conf"),
				"locals": cty.MapVal(map[string]cty.Value{
					"temp_variable": cty.StringVal("fake me"),
				}),
			})

			_, err := templateCommand(fakeContext)

			expect.NoErrors(t, err)

			command := fakeContext.CommandRan()
			expect.EqualFatal(t, len(command.Command), 4)
			assert.Equal(t, command.Command[0], `echo "c245841ee91cee0ea02879e88783c9269f508f83 $HOME/file.conf" | sha1sum -c -`)
			assert.Equal(t, command.Command[1], `if [ $? = 0 ]; then exit 40; fi`)
			assert.Equal(t, command.Command[2], "filebase64encoded=VGhpcyBpcyBhbiBleGFtcGxlIHRlbXBsYXRlIGZpbGUKCnRlbXBsYXRlX2ZpbGUuY29uZgo=")
			assert.Equal(t, command.Command[3], "echo $filebase64encoded | base64 -d > $HOME/file.conf")
		})
	})
}
