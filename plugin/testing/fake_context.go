package testing

import (
	"github.com/pythonandchips/azad/plugin"
	"github.com/zclconf/go-cty/cty"
)

// FakeContext for testing plugins
type FakeContext struct {
	vars       map[string]cty.Value
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
func (context *FakeContext) SetVars(vars map[string]cty.Value) {
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
	return context.vars[key].AsString()
}

// Exists check the key exists
func (context FakeContext) Exists(key string) bool {
	_, ok := context.vars[key]
	return ok
}

// GetMap return value as a map
func (context FakeContext) GetMap(key string) map[string]string {
	out := map[string]string{}
	v, ok := context.vars[key]
	if !ok {
		return out
	}
	val := v.AsValueMap()
	for k, v := range val {
		out[k] = v.AsString()
	}
	return out
}

// GetWithDefault return value or supplied default
func (context FakeContext) GetWithDefault(key, def string) string {
	val, ok := context.vars[key]
	if !ok {
		return def
	}
	return val.AsString()
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
	v, ok := context.vars[key]
	if !ok {
		return false
	}
	return v.True()
}
