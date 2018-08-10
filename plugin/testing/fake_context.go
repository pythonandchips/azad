package testing

import (
	"github.com/pythonandchips/azad/plugin"
)

// FakeContext for testing plugins
type FakeContext struct {
	vars       map[string]string
	env        map[string]string
	user       string
	stdout     string
	stderr     string
	rolePath   string
	rootPath   string
	ranCommand plugin.Command
	runErr     error
}

// NewFakeContext use as a mock for real context in tests
func NewFakeContext() *FakeContext {
	return &FakeContext{}
}

// CommandRan returns the command passed to the run function of the
// context
func (context *FakeContext) CommandRan() plugin.Command {
	return context.ranCommand
}

// SetVars for the fake context
func (context *FakeContext) SetVars(vars map[string]string) {
	context.vars = vars
}

// SetRolePath for the fake context
func (context *FakeContext) SetRolePath(rolePath string) {
	context.rolePath = rolePath
}

// Run command against a fake connection
func (context *FakeContext) Run(command plugin.Command) error {
	context.ranCommand = command
	return context.runErr
}

// Stdout retrieve the result of stdout sent by the last run
func (context FakeContext) Stdout() string {
	return context.stdout
}

// User to run command with
func (context FakeContext) User() string {
	return context.user
}

// Stderr retrieve the result of stdout sent by the last run
func (context FakeContext) Stderr() string {
	return context.stderr
}

// Get the configuration value for the task
func (context FakeContext) Get(key string) string {
	return context.vars[key]
}

// Exists check the key exists
func (context FakeContext) Exists(key string) bool {
	_, ok := context.vars[key]
	return ok
}

// Variables raw variables for context
func (context FakeContext) Variables() map[string]string {
	return context.vars
}

// GetWithDefault return value or supplied default
func (context FakeContext) GetWithDefault(key, def string) string {
	val, ok := context.vars[key]
	if !ok {
		return def
	}
	return val
}

// PlaybookPath absolute path to root of running playbook
func (context FakeContext) PlaybookPath() string {
	return context.rootPath
}

// RolePath absolute path to root of running role
func (context FakeContext) RolePath() string {
	return context.rolePath
}

// IsTrue checks if key exists and matches "true"
func (context FakeContext) IsTrue(key string) bool {
	return context.Get(key) == "true"
}
