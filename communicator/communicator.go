package communicator

import (
	"fmt"

	"github.com/pythonandchips/azad/communicator/awsinventory"
	"github.com/pythonandchips/azad/communicator/core"
	"github.com/pythonandchips/azad/plugin"
)

var schemas = map[string]plugin.Schema{}

var getPlugin = func(pluginName string) (plugin.Schema, error) {
	if _, ok := schemas[pluginName]; !ok {
		switch pluginName {
		case core.Name:
			schemas[pluginName] = core.GetSchema()
		case awsinventory.Name:
			schemas[pluginName] = awsinventory.GetSchema()
		default:
			return plugin.Schema{}, fmt.Errorf("plugin %s not found", pluginName)
		}
	}
	return schemas[pluginName], nil
}

// GetTask returns a task for the specified plugin e.g. core.bash
func GetTask(pluginName, taskName string) (plugin.Task, error) {
	schema, err := getPlugin(pluginName)
	if err != nil {
		return plugin.Task{}, err
	}
	task, ok := schema.Tasks[taskName]
	if !ok {
		return plugin.Task{}, fmt.Errorf("plugin %s does not contain task %s", pluginName, taskName)
	}
	return task, nil
}

// GetInventory returns the inventory for the specified plugin and service e.g. aws.ec2
func GetInventory(pluginName, serviceName string) (plugin.Inventory, error) {
	schema, err := getPlugin(pluginName)
	if err != nil {
		return plugin.Inventory{}, err
	}
	service, ok := schema.Inventory[serviceName]
	if !ok {
		return plugin.Inventory{}, fmt.Errorf("plugin %s does not contain inventory %s", pluginName, serviceName)
	}
	return service, nil
}
