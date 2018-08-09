package aptget

import (
	"fmt"
	"path"

	"github.com/pythonandchips/azad/plugin"
)

func installConfig() plugin.Task {
	return plugin.Task{
		Fields: []plugin.Field{
			{Name: "package", Type: "String", Required: true},
			{Name: "update", Type: "String"},
			{Name: "deb", Type: "String"},
		},
		Run: installCommand,
	}
}

func installCommand(context plugin.Context) (map[string]string, error) {
	commands := []string{}
	if context.Exists("deb") {
		url := context.Get("deb")
		basename := path.Base(url)
		commands = append(commands, fmt.Sprintf(`wget -qO /tmp/%s %s`, basename, url))
		commands = append(commands, fmt.Sprintf(`dpkg -i /tmp/%s`, basename))
		commands = append(commands, fmt.Sprintf(`rm /tmp/%s`, basename))
		if !context.IsTrue("update") {
			commands = append(commands, "apt-get update -yq")
		}
	}
	if context.IsTrue("update") {
		commands = append(commands, "apt-get update -yq")
	}
	commands = append(
		commands,
		fmt.Sprintf("apt-get install -yq %s", context.Get("package")),
	)
	command := plugin.Command{
		Command: commands,
	}
	err := context.Run(command)
	return map[string]string{}, err
}
