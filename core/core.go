package core

import (
	"fmt"
	"plugin"

	"github.com/pythonandchips/azad/schema"
)

var corePlugin CorePlugin

func New() CorePlugin {
	if len(corePlugin.funcs) == 0 {
		corePlugin = CorePlugin{
			funcs: map[string]interface{}{
				"Schema": getSchema,
			},
		}
	}
	return corePlugin
}

type CorePlugin struct {
	funcs map[string]interface{}
}

func (corePlugin CorePlugin) Lookup(method string) (plugin.Symbol, error) {
	if funcs, ok := corePlugin.funcs[method]; ok {
		return funcs, nil
	}
	return nil, fmt.Errorf("%s not available in core", method)
}

func getSchema() schema.Schema {
	return schema.Schema{
		Tasks: []schema.Task{
			{
				Name: "bash", Fields: []schema.Field{
					{Name: "command", Type: "String", Required: true},
					{Name: "chdir", Type: "String", Required: false},
				},
				Run: runCommand,
			},
		},
	}
}

func runCommand(context schema.Context) error {
	command := schema.Command{
		Interpreter: "bash",
		Command: []string{
			context.Get("command"),
		},
	}
	err := context.Run(command)
	return err
}
