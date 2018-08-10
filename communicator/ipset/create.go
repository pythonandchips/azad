package ipset

import (
	"fmt"

	"github.com/pythonandchips/azad/plugin"
)

func createConfig() plugin.Task {
	return plugin.Task{
		Fields: []plugin.Field{
			{Name: "set", Type: "String", Required: true},
			{Name: "entry", Type: "String", Required: true},
			{Name: "dest", Type: "String", Required: false},
		},
		Run: createCommand,
	}
}
func createCommand(context plugin.Context) (map[string]string, error) {
	commands := []string{
		fmt.Sprintf("ipset create %s %s", context.Get("set"), context.Get("entry")),
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
