package aptget

import (
	"fmt"

	"github.com/pythonandchips/azad/plugin"
)

func updateConfig() plugin.Task {
	return plugin.Task{
		Fields: []plugin.Field{},
		Run:    update,
	}
}

func update(context plugin.Context) (map[string]string, error) {
	command := plugin.Command{
		Interpreter: "sh",
		Command: []string{
			fmt.Sprintf("apt-get update -yq"),
		},
	}
	err := context.Run(command)
	return map[string]string{}, err
}
