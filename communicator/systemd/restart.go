package systemd

import (
	"fmt"

	"github.com/pythonandchips/azad/plugin"
)

func restartConfig() plugin.Task {
	return plugin.Task{
		Fields: []plugin.Field{
			{Name: "service", Type: "String", Required: true},
		},
		Run: restart,
	}
}

// Restart a systemd process
func restart(context plugin.Context) (map[string]string, error) {
	command := plugin.Command{
		Interpreter: "sh",
		Command: []string{
			fmt.Sprintf("systemd restart %s", context.Get("service")),
		},
	}
	err := context.Run(command)
	return map[string]string{}, err
}
