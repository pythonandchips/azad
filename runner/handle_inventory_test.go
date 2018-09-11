package runner

import (
	"fmt"
	"testing"

	"github.com/hashicorp/hcl2/hcl"
	"github.com/pythonandchips/azad/plugin"
	"github.com/pythonandchips/azad/steps"
	"github.com/stretchr/testify/assert"
)

func TestHandleInventory(t *testing.T) {
	t.Run("fetches servers from remote", func(t *testing.T) {
		inventoryStep, store := testHandleInventorySetup()
		err := handleInventory(inventoryStep, &store)
		if err != nil {
			t.Fatalf("unexpected err: %s", err)
		}
		assert.Equal(t, len(store.servers), 2)
	})
	t.Run("when plugin cannot be found", func(t *testing.T) {
		inventoryStep, store := testHandleInventorySetup()
		fakeGetInventoryError()
		err := handleInventory(inventoryStep, &store)
		assert.Equal(t, "plugin not found", err.Error())
	})
	t.Run("when an issue is raised with the body content", func(t *testing.T) {
		inventoryStep, store := testHandleInventorySetup()
		inventoryStep.Body = &TestBody{
			err: hcl.Diagnostics{
				{Severity: hcl.DiagInvalid, Summary: "missing attrs"},
			},
		}
		err := handleInventory(inventoryStep, &store)
		if err == nil {
			t.Fatalf("expected error but got none")
		}
		assert.Equal(t, "<nil>: missing attrs; ", err.Error())
	})
	t.Run("when unable to evaluate variable", func(t *testing.T) {
		inventoryStep, store := testHandleInventorySetup()
		inventoryStep.Body = &TestBody{
			attributes: map[string]*hcl.Attribute{
				"access_id":  testExpression("access_id", `rank`),
				"secret_key": testExpression("secret_key", `"ABCDEF09876543"`),
			},
		}
		err := handleInventory(inventoryStep, &store)
		if err == nil {
			t.Fatalf("expected error but got none")
		}
		assert.Equal(t, "1 error occurred:\n\t* variable not found: rank\n\n", err.Error())
	})
	t.Run("when inventory raises an error", func(t *testing.T) {
		inventoryStep, store := testHandleInventorySetup()
		fakeCallToInventoryError()
		err := handleInventory(inventoryStep, &store)
		if err == nil {
			t.Fatalf("expected error but got none")
		}
		assert.Equal(t, "error running inventory", err.Error())
	})
}

func testHandleInventorySetup() (steps.InventoryStep, store) {
	fakeGetInventory()
	store := store{
		variables: variables{},
		servers:   servers{},
	}
	inventoryStep := steps.InventoryStep{
		Type: "aws.ec2",
		Body: &TestBody{
			attributes: map[string]*hcl.Attribute{
				"access_id":  testExpression("access_id", `"abcdef123456"`),
				"secret_key": testExpression("secret_key", `"ABCDEF09876543"`),
			},
		},
	}
	return inventoryStep, store
}

func fakeGetInventory() {
	getInventory = func(pluginName, serviceName string) (plugin.Inventory, error) {
		return plugin.Inventory{
			Fields: []plugin.Field{
				{Name: "access_id", Required: true},
				{Name: "secret_key"},
			},
			Run: func(plugin.InventoryContext) ([]plugin.Resource, error) {
				return []plugin.Resource{
					{ConnectOn: "10.0.0.1", Groups: []string{"development", "kibana_server"}},
					{ConnectOn: "10.0.0.2", Groups: []string{"kibana_server"}},
				}, nil
			},
		}, nil
	}
}

func fakeCallToInventoryError() {
	getInventory = func(pluginName, serviceName string) (plugin.Inventory, error) {
		return plugin.Inventory{
			Fields: []plugin.Field{
				{Name: "access_id", Required: true},
				{Name: "secret_key"},
			},
			Run: func(plugin.InventoryContext) ([]plugin.Resource, error) {
				return []plugin.Resource{}, fmt.Errorf("error running inventory")
			},
		}, nil
	}
}

func fakeGetInventoryError() {
	getInventory = func(pluginName, serviceName string) (plugin.Inventory, error) {
		return plugin.Inventory{}, fmt.Errorf("plugin not found")
	}
}
