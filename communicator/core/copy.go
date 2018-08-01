package core

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/pythonandchips/azad/plugin"
)

func copyConfig() plugin.Task {
	return plugin.Task{
		Fields: []plugin.Field{
			{Name: "source", Type: "String", Required: true},
			{Name: "dest", Type: "String", Required: true},
			{Name: "mode", Type: "String", Required: false},
		},
		Run: copyCommand,
	}
}

func copyCommand(context plugin.Context) (map[string]string, error) {
	templatePath := filepath.Join(context.RolePath(), "files", context.Get("source"))
	data, _ := ioutil.ReadFile(templatePath)
	commands := []string{
		"filecont=<<- HEREDOC",
		string(data),
		"HEREDOC",
		fmt.Sprintf("echo $filecount > %s", context.Get("dest")),
	}
	if context.Exists("mode") {
		commands = append(
			commands,
			fmt.Sprintf("chmod %s %s", context.Get("mode"), context.Get("dest")),
		)
	}
	command := plugin.Command{
		Interpreter: "sh",
		Command:     commands,
	}
	err := context.Run(command)
	return map[string]string{}, err
}
