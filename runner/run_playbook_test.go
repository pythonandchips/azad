package runner

import (
	"testing"

	"github.com/pythonandchips/azad/conn"
	"github.com/pythonandchips/azad/logger"
	"github.com/pythonandchips/azad/plugin"
	"github.com/pythonandchips/azad/steps"
)

func TestRunPlaybook(t *testing.T) {
	logger.StubLogger()
	conn.Simulate()
	playbookSteps := steps.PlaybookSteps{
		StepList: steps.StepList{
			defaultServerStep(),
			defaultInventoryStep(),
			defaultVariableStep(),
			defaultInputTask(),
			mapVariableStep(),
			arrayVariableStep(),
			defaultContext(),
		},
	}
	playbookPath := "/home/user/playbooks/basic.az"
	parsePlaybook = func(playbookFilePath string) (steps.PlaybookSteps, error) {
		return playbookSteps, nil
	}
	roleList = []steps.RoleContainer{
		defaultRoleContainer(),
	}
	getInventory = func(string, string) (plugin.Inventory, error) {
		return plugin.Inventory{
			Fields: []plugin.Field{
				{Name: "access_id"},
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
	getTask = func(pluginName, taskName string) (plugin.Task, error) {
		return testPluginTask(), nil
	}
	getInput = func(pluginName, taskName string) (plugin.Input, error) {
		return testPluginInput(), nil
	}

	err := RunPlaybook(playbookPath, Config{})

	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
}
