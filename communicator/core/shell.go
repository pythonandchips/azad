package core

import "github.com/pythonandchips/azad/plugin"

func shellConfig() plugin.Task {
	return plugin.Task{
		Fields: []plugin.Field{
			{Name: "command", Type: "String", Required: true},
		},
		Run: shell,
	}
}

func shell(context plugin.Context) (map[string]string, error) {
	command := plugin.Command{
		Interpreter: "sh",
		Command: []string{
			context.Get("command"),
		},
	}
	err := context.Run(command)
	return map[string]string{}, err
}
