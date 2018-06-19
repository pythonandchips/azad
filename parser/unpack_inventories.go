package parser

import (
	"github.com/hashicorp/hcl2/hcl"
	"github.com/pythonandchips/azad/azad"
)

func unpackInventories(inventoryDescriptions []inventoryDescription, evalContext *hcl.EvalContext) ([]azad.Inventory, error) {
	inventories := []azad.Inventory{}
	for _, inventoryDescription := range inventoryDescriptions {
		inventory := azad.Inventory{
			Name: inventoryDescription.Name,
		}
		attributes := map[string]string{}
		attributesList, _ := inventoryDescription.Config.JustAttributes()

		for _, attr := range attributesList {
			value, _ := attr.Expr.Value(evalContext)
			attributes[attr.Name] = value.AsString()
		}
		inventory.Attributes = attributes
		inventories = append(inventories, inventory)
	}
	return inventories, nil
}
