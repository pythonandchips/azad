package core

import "github.com/pythonandchips/azad/plugin"

func bashConfig() plugin.Task {
	return plugin.Task{
		Fields: []plugin.Field{
			{Name: "command", Type: "String", Required: true},
			{Name: "chdir", Type: "String", Required: false},
		},
		Run: bash,
	}
}

func bash(context plugin.Context) (map[string]string, error) {
	command := plugin.Command{
		Interpreter: "bash",
		Command: []string{
			context.Get("command"),
		},
	}
	err := context.Run(command)
	return map[string]string{}, err
}
