package conn

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommandGenerateFile(t *testing.T) {
	command := Command{
		Interpreter: "sh",
		Command: []string{
			"ls -al",
			"touch /tmp/azad.run",
		},
		Env: map[string]string{
			"HELLO": "WORLD",
		},
	}
	t.Run("generate a file to run on server", func(t *testing.T) {
		expectedFile := `#!/usr/bin/env sh

HELLO='WORLD'

ls -al
touch /tmp/azad.run`
		assert.Equal(t, command.generateFile(), expectedFile)
	})
}
