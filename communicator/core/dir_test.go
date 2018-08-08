package core

import (
	"testing"

	"github.com/pythonandchips/azad/expect"
	plugintesting "github.com/pythonandchips/azad/plugin/testing"
	"github.com/stretchr/testify/assert"
)

func TestDir(t *testing.T) {
	t.Run("when command is successful", func(t *testing.T) {
		t.Run("with only path specified", func(t *testing.T) {
			fakeContext := plugintesting.NewFakeContext()
			fakeContext.SetVars(map[string]string{
				"path": "/home/he-man",
			})
			_, err := dir(fakeContext)
			expect.NoErrors(t, err)

			command := fakeContext.CommandRan()
			expect.EqualFatal(t, len(command.Command), 2)
			assert.Equal(t, command.Command[0], "if [ -d /home/he-man ]; then exit 0; fi")
			assert.Equal(t, command.Command[1], "mkdir -p /home/he-man")
		})

		t.Run("with path and owner specified", func(t *testing.T) {
			fakeContext := plugintesting.NewFakeContext()
			fakeContext.SetVars(map[string]string{
				"path":  "/home/he-man",
				"owner": "heman",
			})
			_, err := dir(fakeContext)
			expect.NoErrors(t, err)

			command := fakeContext.CommandRan()
			expect.EqualFatal(t, len(command.Command), 3)
			assert.Equal(t, command.Command[0], "if [ -d /home/he-man ]; then exit 0; fi")
			assert.Equal(t, command.Command[1], "mkdir -p /home/he-man")
			assert.Equal(t, command.Command[2], "chown heman:heman /home/he-man")
		})

		t.Run("with path, owner, group specified", func(t *testing.T) {
			fakeContext := plugintesting.NewFakeContext()
			fakeContext.SetVars(map[string]string{
				"path":  "/home/he-man",
				"owner": "heman",
				"group": "grayskull",
			})
			_, err := dir(fakeContext)
			expect.NoErrors(t, err)

			command := fakeContext.CommandRan()
			expect.EqualFatal(t, len(command.Command), 3)
			assert.Equal(t, command.Command[0], "if [ -d /home/he-man ]; then exit 0; fi")
			assert.Equal(t, command.Command[1], "mkdir -p /home/he-man")
			assert.Equal(t, command.Command[2], "chown heman:grayskull /home/he-man")
		})
	})
}
