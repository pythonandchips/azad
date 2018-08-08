package aptget

import (
	"fmt"

	"github.com/pythonandchips/azad/plugin"
)

func installConfig() plugin.Task {
	return plugin.Task{
		Fields: []plugin.Field{
			{Name: "package", Type: "String", Required: true},
		},
		Run: installCommand,
	}
}

func installCommand(context plugin.Context) (map[string]string, error) {
	command := plugin.Command{
		Command: []string{
			fmt.Sprintf("apt-get install -yq %s", context.Get("package")),
		},
	}
	err := context.Run(command)
	return map[string]string{}, err
}
