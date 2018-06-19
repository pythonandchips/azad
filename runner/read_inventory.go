package runner

import (
	"github.com/pythonandchips/azad/azad"
	"github.com/pythonandchips/azad/communicator"
	"github.com/pythonandchips/azad/plugin"
)

func readInventory(playbook azad.Playbook) (azad.Playbook, error) {
	for _, inventory := range playbook.Inventories {
		pluginName, err := inventory.PluginName()
		if err != nil {
			return playbook, err
		}
		serviceName, err := inventory.ServiceName()
		if err != nil {
			return playbook, err
		}
		inventorySchema, _ := communicator.GetInventory(pluginName, serviceName)
		inventoryContext := plugin.NewInventoryContext(inventory.Attributes)
		resources, _ := inventorySchema.Run(inventoryContext)
		for _, resource := range resources {
			for _, group := range resource.Groups {
				playbook.AddAddressToServerByGroup(group, resource.ConnectOn)
			}
		}
	}
	return playbook, nil
}
