package core

import "github.com/pythonandchips/azad/plugin"

func templateConfig() plugin.Task {
	return plugin.Task{
		Fields: []plugin.Field{
			{Name: "source", Type: "String", Required: true},
			{Name: "dest", Type: "String", Required: true},
		},
		Run: template,
	}
}

func template(context plugin.Context) (map[string]string, error) {
	return map[string]string{}, nil
}
