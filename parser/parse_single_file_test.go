package parser

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/pythonandchips/azad/expect"
	"github.com/pythonandchips/azad/steps"
	"github.com/stretchr/testify/assert"
)

func TestParseSingleFilePlaybook(t *testing.T) {
	wd, _ := os.Getwd()
	filePath := filepath.Join(
		wd, "fixtures", "basic.az",
	)
	playbookSteps, err := PlaybookSteps(filePath)

	stepList := playbookSteps.StepList

	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	expect.EqualFatal(t, len(stepList), 5)

	t.Run("server step", func(t *testing.T) {
		serverStep := stepList[0].(steps.ServerStep)
		assert.Equal(t, serverStep.Name, "kibana_server")
		assert.NotNil(t, serverStep.Addresses)
	})

	t.Run("inventory step", func(t *testing.T) {
		inventoryStep := stepList[1].(steps.InventoryStep)
		assert.Equal(t, inventoryStep.Type, "aws.ec2")
		assert.NotNil(t, inventoryStep.Body)
	})

	t.Run("variable step", func(t *testing.T) {
		variableStep := stepList[2].(steps.VariableStep)
		assert.Equal(t, variableStep.Name, "base_path")
		assert.NotNil(t, variableStep.Type)
		assert.NotNil(t, variableStep.Value)
	})

	t.Run("input step", func(t *testing.T) {
		inputStep := stepList[3].(steps.InputStep)
		assert.Equal(t, inputStep.Type, "credstash.get")
		assert.NotNil(t, inputStep.Body)
	})

	t.Run("context step", func(t *testing.T) {
		contextStep := stepList[4].(steps.ContextContainer)
		assert.Equal(t, contextStep.Name, "kibana server")
		assert.NotNil(t, contextStep.User)
		assert.NotNil(t, contextStep.ApplyTo)
		assert.Equal(t, len(contextStep.Steps), 4)

		t.Run("include step", func(t *testing.T) {
			includeStep := contextStep.Steps[3].(steps.IncludesStep)
			assert.NotNil(t, includeStep.Roles)
		})
	})

	t.Run("role container", func(t *testing.T) {
		roleStep := playbookSteps.RoleList[0]
		assert.Equal(t, roleStep.Name, "elasticsearch")
		assert.NotNil(t, roleStep.User)
		assert.Equal(t, len(roleStep.Steps), 4)
	})
}
