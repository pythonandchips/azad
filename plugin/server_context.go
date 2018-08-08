package plugin

import "github.com/pythonandchips/azad/conn"

// ServerContext for task run. Main representation of data for task run
type serverContext struct {
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
func (context *serverContext) Run(command Command) error {
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

func (context serverContext) User() string {
	return context.user
}

// Stdout retrieve the result of stdout sent by the last run
func (context serverContext) Stdout() string {
	return context.stdout
}

// Stderr retrieve the result of stdout sent by the last run
func (context serverContext) Stderr() string {
	return context.stderr
}

// Get the configuration value for the task
func (context serverContext) Get(key string) string {
	return context.vars[key]
}

// Exists check the key exists
func (context serverContext) Exists(key string) bool {
	_, ok := context.vars[key]
	return ok
}

// Variables raw variables for context
func (context serverContext) Variables() map[string]string {
	return context.vars
}

// GetWithDefault return value or supplied default
func (context serverContext) GetWithDefault(key, def string) string {
	val, ok := context.vars[key]
	if !ok {
		return def
	}
	return val
}

// PlaybookPath absolute path to root of running playbook
func (context serverContext) PlaybookPath() string {
	return context.rootPath
}

// RolePath absolute path to root of running role
func (context serverContext) RolePath() string {
	return context.rolePath
}
