package plugin

import (
	"github.com/pythonandchips/azad/conn"
)

// NewContext create a new context for passing to a plugin to run a task
func NewContext(vars map[string]string, conn conn.Conn, user, rootPath, rolePath string) Context {
	return &serverContext{
		vars:     vars,
		conn:     conn,
		user:     user,
		rootPath: rootPath,
		rolePath: rolePath,
	}
}

// NewInventoryContext create a new context for passing to a plugin to run a task
func NewInventoryContext(vars map[string]string) Context {
	return &serverContext{
		vars: vars,
	}
}

// Context for task run. Main representation of data for task run
type Context interface {
	Run(Command) error
	Stdout() string
	Stderr() string
	User() string
	Get(string) string
	Exists(string) bool
	Variables() map[string]string
	GetWithDefault(string, string) string
	PlaybookPath() string
	RolePath() string
}
