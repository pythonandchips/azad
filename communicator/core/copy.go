package core

import "github.com/pythonandchips/azad/plugin"

func copyConfig() plugin.Task {
	return plugin.Task{
		Fields: []plugin.Field{
			{Name: "source", Type: "String", Required: true},
			{Name: "dest", Type: "String", Required: true},
			{Name: "mode", Type: "String", Required: false},
		},
		Run: copyCommand,
	}
}

func copyCommand(context plugin.Context) (map[string]string, error) {
	return map[string]string{}, nil
}
