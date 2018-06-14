package plugins

import (
	"fmt"
	"plugin"

	"github.com/pythonandchips/azad/schema"
)

var pluginList = map[string]Plugin{}

// Plugin plugin
type Plugin interface {
	TaskExists(string) error
	GetTask(string) (schema.Task, error)
}

// IPlugin blah
type IPlugin interface {
	Lookup(string) (plugin.Symbol, error)
}

type goPlugin struct {
	plugin IPlugin
	schema schema.Schema
	name   string
}

func (p goPlugin) TaskExists(taskName string) error {
	for _, task := range p.schema.Tasks {
		if task.Name == taskName {
			return nil
		}
	}
	return fmt.Errorf("Task %s not found in plugin %s", taskName, p.name)
}

// GetTask task
func (p goPlugin) GetTask(taskName string) (schema.Task, error) {
	for _, task := range p.schema.Tasks {
		if task.Name == taskName {
			return task, nil
		}
	}
	return schema.Task{}, fmt.Errorf("Task %s not found in plugin %s", taskName, p.name)
}
