package ipset

import (
	"fmt"

	"github.com/pythonandchips/azad/plugin"
)

func openPortConfig() plugin.Task {
	return plugin.Task{
		Fields: []plugin.Field{
			{Name: "set", Type: "String", Required: true},
			{Name: "port", Type: "String", Required: true},
			{Name: "dest", Type: "String"},
		},
		Run: openPortCommand,
	}
}

func openPortCommand(context plugin.Context) (map[string]string, error) {
	commands := []string{
		fmt.Sprintf(`ipset -! add %s %s`, context.Get("set"), context.Get("port")),
	}
	if context.Exists("dest") {
		commands = append(
			commands,
			fmt.Sprintf(`ipset save %s > %s`, context.Get("set"), context.Get("dest")),
		)
	}
	command := plugin.Command{
		Command: commands,
	}
	err := context.Run(command)
	return map[string]string{}, err
}
