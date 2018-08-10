package core

import (
	"fmt"

	"github.com/pythonandchips/azad/plugin"
)

func removeFileConfig() plugin.Task {
	return plugin.Task{
		Fields: []plugin.Field{
			{Name: "path", Type: "String", Required: true},
		},
		Run: removeFileCommand,
	}
}

func removeFileCommand(context plugin.Context) (map[string]string, error) {
	command := plugin.Command{
		Command: []string{
			fmt.Sprintf("rm %s", context.Get("path")),
		},
	}
	err := context.Run(command)
	return map[string]string{}, err
}
