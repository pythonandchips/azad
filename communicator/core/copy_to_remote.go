package core

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/pythonandchips/azad/plugin"
	"github.com/pythonandchips/azad/plugin/helpers"
)

func copyToRemoteConfig() plugin.Task {
	return plugin.Task{
		Fields: []plugin.Field{
			{Name: "source", Type: "String", Required: true},
			{Name: "dest", Type: "String", Required: true},
			{Name: "mode", Type: "String", Required: false},
		},
		Run: copyToRemoteCommand,
	}
}

func copyToRemoteCommand(context plugin.Context) (map[string]string, error) {
	templatePath := filepath.Join(context.RolePath(), "files", context.Get("source"))
	data, err := ioutil.ReadFile(templatePath)
	if err != nil {
		return map[string]string{}, err
	}
	commands := []string{}
	commands = append(commands, helpers.Checksum(data, context.Get("dest"))...)
	commands = append(commands, helpers.WriteEncodedFile(data, context.Get("dest"))...)

	if context.Exists("mode") {
		commands = append(
			commands,
			fmt.Sprintf("chmod %s %s", context.Get("mode"), context.Get("dest")),
		)
	}
	command := plugin.Command{
		Command: commands,
	}
	err = context.Run(command)
	return map[string]string{}, err
}
