package core

import (
	"fmt"
	"plugin"

	"github.com/pythonandchips/azad/conn"
	"github.com/pythonandchips/azad/logger"
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

func runCommand(vars map[string]string, connection conn.Conn) error {
	command := conn.Command{
		Interpreter: "bash",
		Command: []string{
			vars["command"],
		},
	}

	commandResponse, err := connection.Run(command)
	logger.Debug(commandResponse.Stdout())
	logger.Debug(commandResponse.Stderr())
	return err
}
