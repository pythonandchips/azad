package core

import (
	"fmt"

	"github.com/pythonandchips/azad/plugin"
)

func dirConfig() plugin.Task {
	return plugin.Task{
		Fields: []plugin.Field{
			{Name: "path", Type: "String", Required: true},
			{Name: "owner", Type: "String", Required: false},
			{Name: "group", Type: "String", Required: false},
		},
		Run: dir,
	}
}

func dir(context plugin.Context) (map[string]string, error) {
	commands := []string{
		checkDirExists(context.Get("path")),
		fmt.Sprintf(`mkdir -p %s`, context.Get("path")),
	}
	if context.Exists("owner") {
		owner := context.Get("owner")
		commands = append(commands,
			fmt.Sprintf("chown %s:%s %s",
				owner,
				context.GetWithDefault("group", owner),
				context.Get("path"),
			),
		)
	}
	command := plugin.Command{
		Command: commands,
	}
	err := context.Run(command)
	return map[string]string{}, err
}

func checkDirExists(path string) string {
	return fmt.Sprintf("if [ -d %s ]; then exit(0); fi", path)
}
