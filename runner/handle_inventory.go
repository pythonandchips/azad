package runner

import (
	"github.com/hashicorp/hcl2/hcl"
	"github.com/pythonandchips/azad/communicator"
	"github.com/pythonandchips/azad/plugin"
	"github.com/pythonandchips/azad/steps"
	"github.com/zclconf/go-cty/cty"
)

func handleInventory(inventoryStep steps.InventoryStep, store serverStore) error {
	inventorySchema, err := inventorySchema(inventoryStep)
	if err != nil {
		return err
	}
	attributes, err := attributesForSchema(inventoryStep.Body, inventorySchema.Fields, store)
	if err != nil {
		return err
	}
	inventoryContext := plugin.NewInventoryContext(attributes)
	resources, err := inventorySchema.Run(inventoryContext)
	if err != nil {
		return err
	}
	servers := servers{}
	for _, resource := range resources {
		server := server{
			address: resource.ConnectOn,
			group:   resource.Groups,
		}
		servers = append(servers, server)
	}
	store.addServers(servers)
	return nil
}

func inventorySchema(inventoryStep steps.InventoryStep) (plugin.Inventory, error) {
	pluginName, err := inventoryStep.PluginName()
	if err != nil {
		return plugin.Inventory{}, err
	}
	serviceName, err := inventoryStep.ServiceName()
	if err != nil {
		return plugin.Inventory{}, err
	}
	return getInventory(pluginName, serviceName)
}

func attributesForSchema(body hcl.Body, fields []plugin.Field, store evalStore) (map[string]cty.Value, error) {
	attributes := map[string]cty.Value{}
	inventoryBodySchema := createSchemaFromFields(fields)
	content, contentErr := body.Content(inventoryBodySchema)
	if contentErr != nil {
		return attributes, contentErr
	}
	for name, attr := range content.Attributes {
		val, err := store.evalVariable(attr, allowString)
		if err != nil {
			return attributes, err
		}
		attributes[name] = val
	}
	return attributes, nil
}

var getInventory = func(pluginName, serviceName string) (plugin.Inventory, error) {
	return communicator.GetInventory(pluginName, serviceName)
}
