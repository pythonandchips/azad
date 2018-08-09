package ipset

import "github.com/pythonandchips/azad/plugin"

// Name of plugin as used in playbook files
//
// Core is the default plugin and the plugin name can be omitted.
// e.g. `core.bash` and `bash` will run the same task
const Name = "ipset"

// GetSchema returns the schema for the ipset plugin
func GetSchema() plugin.Schema {
	return plugin.Schema{
		Tasks: map[string]plugin.Task{
			"create":    createConfig(),
			"open-port": openPortConfig(),
			"save":      saveConfig(),
		},
	}
}
