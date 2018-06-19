package core

import "github.com/pythonandchips/azad/plugin"

// Name of plugin as used in playbook files
//
// Core is the default plugin and the plugin name can be omitted.
// e.g. `core.bash` and `bash` will run the same task
const Name = "core"

// GetSchema returns the schema for the current plugin
//
// Tasks
// 	 bash:
//     command (required: true): the bash command to run on the system e.g. `whois`
//     chdir   (required: fasle): directory to run the command in
func GetSchema() plugin.Schema {
	return plugin.Schema{
		Tasks: map[string]plugin.Task{
			"bash": {
				Fields: []plugin.Field{
					{Name: "command", Type: "String", Required: true},
					{Name: "chdir", Type: "String", Required: false},
				},
				Run: bash,
			},
		},
	}
}

func bash(context plugin.Context) error {
	command := plugin.Command{
		Interpreter: "bash",
		Command: []string{
			context.Get("command"),
		},
	}
	return context.Run(command)
}
