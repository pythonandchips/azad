package ipset

import (
	"fmt"

	"github.com/pythonandchips/azad/plugin"
)

func saveConfig() plugin.Task {
	return plugin.Task{
		Fields: []plugin.Field{
			{Name: "set", Type: "String", Required: true},
			{Name: "dest", Type: "String", Required: true},
		},
		Run: saveCommand,
	}
}

func saveCommand(context plugin.Context) (map[string]string, error) {
	commands := []string{
		fmt.Sprintf(`ipset save %s > %s`, context.Get("set"), context.Get("dest")),
	}
	command := plugin.Command{
		Command: commands,
	}
	err := context.Run(command)
	return map[string]string{}, err
}
