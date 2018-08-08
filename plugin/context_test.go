package plugin

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zclconf/go-cty/cty"
)

func TestServerContextRun(t *testing.T) {
	command := Command{
		Interpreter: "bash",
		Command:     []string{"pwd"},
	}
	response := FakeResponse{
		stdout: "stdout",
		stderr: "stderr",
	}
	fakeConn := &FakeConn{response: response}
	context := serverContext{
		conn: fakeConn,
		env: map[string]string{
			"env": "dev",
		},
		vars: map[string]cty.Value{
			"variable": cty.StringVal("value"),
		},
		user: "root",
	}

	context.Run(command)
	t.Run("returns stdout form command", func(t *testing.T) {
		assert.Equal(t, context.Stdout(), response.stdout)
	})
	t.Run("returns stderr form command", func(t *testing.T) {
		assert.Equal(t, context.Stderr(), response.stderr)
	})
	t.Run("passes the command to be ran on connection", func(t *testing.T) {
		assert.Equal(t, fakeConn.command.Interpreter, command.Interpreter)
		assert.Equal(t, fakeConn.command.Command, command.Command)
		assert.Equal(t, fakeConn.command.User, context.User())
		assert.Equal(t, fakeConn.command.Env, command.Env)
	})
	t.Run("get variables for task", func(t *testing.T) {
		assert.Equal(t, context.Get("variable"), "value")
	})
}
