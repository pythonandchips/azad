package plugin

import (
	"github.com/pythonandchips/azad/conn"
)

// NewContext create a new context for passing to a plugin to run a task
func NewContext(vars map[string]string, conn conn.Conn, user, rootPath, rolePath string) Context {
	return Context{
		vars:     vars,
		conn:     conn,
		user:     user,
		rootPath: rootPath,
		rolePath: rolePath,
	}
}

// NewInventoryContext create a new context for passing to a plugin to run a task
func NewInventoryContext(vars map[string]string) Context {
	return Context{
		vars: vars,
	}
}

// Context for task run. Main representation of data for task run
type Context struct {
	conn     conn.Conn
	vars     map[string]string
	env      map[string]string
	user     string
	stdout   string
	stderr   string
	rolePath string
	rootPath string
}

// Run command against the ssh connection for the server
func (context *Context) Run(command Command) error {
	cmd := conn.Command{
		Interpreter: command.Interpreter,
		Command:     command.Command,
		User:        context.user,
	}
	response, err := context.conn.Run(cmd)
	context.stdout = response.Stdout()
	context.stderr = response.Stderr()
	return err
}

// User specified to run command
func (context Context) User() string {
	return context.user
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

// Exists check the key exists
func (context Context) Exists(key string) bool {
	_, ok := context.vars[key]
	return ok
}

// GetWithDefault return value or supplied default
func (context Context) GetWithDefault(key, def string) string {
	val, ok := context.vars[key]
	if !ok {
		return def
	}
	return val
}

// PlaybookRoot absolute path to root of running playbook
func (context Context) PlaybookRoot() string {
	return context.rootPath
}

// RoleRoot absolute path to root of running role
func (context Context) RoleRoot() string {
	return context.rolePath
}
