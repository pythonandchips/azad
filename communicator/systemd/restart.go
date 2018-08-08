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
		Run: restartCommand,
	}
}

// Restart a systemd process
func restartCommand(context plugin.Context) (map[string]string, error) {
	command := plugin.Command{
		Command: []string{
			fmt.Sprintf("systemctl restart %s", context.Get("service")),
		},
	}
	err := context.Run(command)
	return map[string]string{}, err
}
