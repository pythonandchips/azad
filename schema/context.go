package schema

import (
	"github.com/pythonandchips/azad/conn"
)

// NewContext create a new context for passing to a plugin to run a task
func NewContext(vars map[string]string, conn conn.Conn) Context {
	return Context{
		vars: vars,
		conn: conn,
	}
}

// Context for task run. Main representation of data for task run
type Context struct {
	conn   conn.Conn
	vars   map[string]string
	env    map[string]string
	stdout string
	stderr string
}

// Run command against the ssh connection for the server
func (context *Context) Run(command Command) error {
	cmd := conn.Command{
		Interpreter: command.Interpreter,
		Command:     command.Command,
		User:        command.User,
	}
	response, err := context.conn.Run(cmd)
	context.stdout = response.Stdout()
	context.stderr = response.Stderr()
	return err
}

// Stdout retrieve the result of stdout sent by the last run
func (context Context) Stdout() string {
	return context.stdout
}

// Stderr retrieve the result of stdout sent by the last run
func (context Context) Stderr() string {
	return context.stderr
}

// Get the configuration value for the task
func (context Context) Get(key string) string {
	return context.vars[key]
}
