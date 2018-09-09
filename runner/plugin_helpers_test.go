package runner

import (
	"fmt"

	"github.com/pythonandchips/azad/plugin"
)

func testPluginTask() plugin.Task {
	return plugin.Task{
		Fields: []plugin.Field{
			{Name: "package", Type: "string", Required: true},
		},
		Run: func(context plugin.Context) (map[string]string, error) {
			command := plugin.Command{
				Command: []string{
					context.Get("package"),
				},
			}
			err := context.Run(command)
			return map[string]string{
				"installed": "true",
			}, err
		},
	}
}

func testFullPluginTask() plugin.Task {
	return plugin.Task{
		Fields: []plugin.Field{
			{Name: "string", Type: "string", Required: true},
			{Name: "variable", Type: "string", Required: true},
			{Name: "interpolation", Type: "string", Required: true},
			{Name: "map-access", Type: "string", Required: true},
			{Name: "array-access", Type: "string", Required: true},
		},
		Run: func(context plugin.Context) (map[string]string, error) {
			command := plugin.Command{
				Command: []string{
					context.Get("string"),
					context.Get("variable"),
					context.Get("interpolation"),
					context.Get("map-access"),
					context.Get("array-access"),
				},
			}
			err := context.Run(command)
			return map[string]string{
				"installed": "true",
			}, err
		},
	}
}

func testErrorPluginTask() plugin.Task {
	return plugin.Task{
		Fields: []plugin.Field{},
		Run: func(context plugin.Context) (map[string]string, error) {
			return map[string]string{}, fmt.Errorf("error running command")
		},
	}
}
