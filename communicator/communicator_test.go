package communicator

import (
	"fmt"
	"testing"

	"github.com/pythonandchips/azad/plugin"
	"github.com/stretchr/testify/assert"
)

func TestGetTask(t *testing.T) {
	aptGetTask := plugin.Task{}
	testSchema := plugin.Schema{
		Tasks: map[string]plugin.Task{
			"apt-get": aptGetTask,
		},
	}
	testPlugin := map[string]plugin.Schema{
		"core": testSchema,
	}
	getPlugin = func(pluginName string) (plugin.Schema, error) {
		if _, ok := testPlugin[pluginName]; !ok {
			return plugin.Schema{}, fmt.Errorf("plugin not found")
		}
		return testPlugin[pluginName], nil
	}

	t.Run("returns the task for the given plugin and service", func(t *testing.T) {
		task, err := GetTask("core", "apt-get")
		if err != nil {
			t.Fatalf("Unexpected Error: %s", err)
		}
		assert.Equal(t, task, aptGetTask)
	})
	t.Run("returns error if plugin does not exists", func(t *testing.T) {
		_, err := GetTask("yum", "install")
		assert.NotNil(t, err, "expected error")
	})
	t.Run("returns error if task does not exists for plugin", func(t *testing.T) {
		_, err := GetTask("core", "systemd")
		assert.NotNil(t, err, "expected error")
	})
}

func TestGetInventory(t *testing.T) {
	ec2Inventory := plugin.Inventory{}
	awsSchema := plugin.Schema{
		Inventory: map[string]plugin.Inventory{
			"ec2": ec2Inventory,
		},
	}
	testPlugin := map[string]plugin.Schema{
		"aws": awsSchema,
	}
	getPlugin = func(pluginName string) (plugin.Schema, error) {
		if _, ok := testPlugin[pluginName]; !ok {
			return plugin.Schema{}, fmt.Errorf("plugin not found")
		}
		return testPlugin[pluginName], nil
	}

	t.Run("returns the inventory for the given plugin and service", func(t *testing.T) {
		inventory, err := GetInventory("aws", "ec2")
		if err != nil {
			t.Fatalf("Unexpected Error: %s", err)
		}
		assert.Equal(t, inventory, ec2Inventory)
	})
	t.Run("returns error if plugin does not exists", func(t *testing.T) {
		_, err := GetInventory("digitalocean", "dropplet")
		assert.NotNil(t, err, "expected error")
	})
	t.Run("returns error if service does not exists for plugin", func(t *testing.T) {
		_, err := GetInventory("aws", "s3")
		assert.NotNil(t, err, "expected error")
	})
}
