package plugin

import (
	"github.com/pythonandchips/azad/conn"
	"github.com/zclconf/go-cty/cty"
)

// Context for task run. Main representation of data for task run
type Context interface {
	Run(Command) error
	Stdout() string
	Stderr() string
	User() string
	Get(string) string
	GetWithDefault(string, string) string
	PlaybookPath() string
	RolePath() string
}

// NewContext create a new context for passing to a plugin to run a task
func NewContext(vars map[string]cty.Value, conn conn.Conn, user, rootPath, rolePath string) Context {
	return &taskContext{
		vars:     vars,
		conn:     conn,
		user:     user,
		rootPath: rootPath,
		rolePath: rolePath,
	}
}

// NewInventoryContext create a new context for passing to a plugin to run a task
func NewInventoryContext(vars map[string]cty.Value) Context {
	return &taskContext{
		vars: vars,
	}
}

// Context for task run. Main representation of data for task run
type taskContext struct {
	conn     conn.Conn
	vars     map[string]cty.Value
	env      map[string]string
	user     string
	stdout   string
	stderr   string
	rolePath string
	rootPath string
}

// Run command against the ssh connection for the server
func (context *taskContext) Run(command Command) error {
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
func (context taskContext) User() string {
	return context.user
}

// Stdout retrieve the result of stdout sent by the last run
func (context taskContext) Stdout() string {
	return context.stdout
}

// Stderr retrieve the result of stdout sent by the last run
func (context taskContext) Stderr() string {
	return context.stderr
}

// Get the configuration value for the task
func (context taskContext) Get(key string) string {
	return context.vars[key].AsString()
}

// Exists check the key exists
func (context Context) Exists(key string) bool {
	_, ok := context.vars[key]
	return ok
}

// GetWithDefault return value or supplied default
func (context taskContext) GetWithDefault(key, def string) string {
	val, ok := context.vars[key]
	if !ok {
		return def
	}
	return val.AsString()
}

// PlaybookRoot absolute path to root of running playbook
func (context taskContext) PlaybookPath() string {
	return context.rootPath
}

// RoleRoot absolute path to root of running role
func (context taskContext) RolePath() string {
	return context.rolePath
}
