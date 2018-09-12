package steps

import (
	"fmt"
	"strings"

	"github.com/hashicorp/hcl2/hcl"
)

// StepList lists steps
type StepList []Step

// Step base for all step types
type Step interface{}

// PlaybookSteps container for steps and roles
type PlaybookSteps struct {
	StepList StepList
	RoleList RoleContainers
}

// ServerStep name and address of a server resource
type ServerStep struct {
	Name      string
	Addresses *hcl.Attribute
}

// InventoryStep fetch server details from an external source
type InventoryStep struct {
	Type string
	Body hcl.Body
}

// PluginName returns the plugin name for the inventory
func (inventoryStep InventoryStep) PluginName() (string, error) {
	parts := strings.Split(inventoryStep.Type, ".")
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid inventory name %s", inventoryStep.Type)
	}
	return parts[0], nil
}

// ServiceName returns the service name for the inventory
func (inventoryStep InventoryStep) ServiceName() (string, error) {
	parts := strings.Split(inventoryStep.Type, ".")
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid inventory name %s", inventoryStep.Type)
	}
	return parts[1], nil
}

// VariableStep config value with name
//
// Available types are string, map, array
type VariableStep struct {
	Name  string
	Type  *hcl.Attribute
	Value *hcl.Attribute
}

// InputStep load variables from an external source
type InputStep struct {
	Type string
	Name string
	Body hcl.Body
}

// PluginName returns the plugin name for the InputStep
func (inputStep InputStep) PluginName() string {
	parts := strings.Split(inputStep.Type, ".")
	if len(parts) == 2 {
		return parts[0]
	}
	return "core"
}

// ServiceName returns the service name for the inventory
func (inputStep InputStep) ServiceName() string {
	parts := strings.Split(inputStep.Type, ".")
	if len(parts) == 2 {
		return parts[1]
	}
	return parts[0]
}

// TaskStep apply an operation to an external resource over ssh
type TaskStep struct {
	Type      string
	Label     string
	User      *hcl.Attribute
	Condition *hcl.Attribute
	Debug     *hcl.Attribute
	WithItems *hcl.Attribute
	Body      hcl.Body
}

// PluginName name of plugin
func (task TaskStep) PluginName() string {
	parts := strings.Split(task.Type, ".")
	if len(parts) == 2 {
		return parts[0]
	}
	return "core"
}

// TaskName name of task
func (task TaskStep) TaskName() string {
	parts := strings.Split(task.Type, ".")
	if len(parts) == 2 {
		return parts[1]
	}
	return parts[0]
}

// IncludesStep brings in other roles and files to make up a full role/context
type IncludesStep struct {
	Roles *hcl.Attribute
}

// ContextContainer connects to a resource and applies all steps to that resource
//
// User - the user that will be used to run the tasks
// ApplyTo - the servers resources to apply the tasks to
type ContextContainer struct {
	Name    string
	Steps   StepList
	User    *hcl.Attribute
	ApplyTo *hcl.Attribute
}

// RoleContainers collection of role containers
type RoleContainers []RoleContainer

// FindByName find a role in container by name.
func (roleContainers RoleContainers) FindByName(name string) (RoleContainer, error) {
	for _, roleContainer := range roleContainers {
		if roleContainer.Name == name {
			return roleContainer, nil
		}
	}
	return RoleContainer{}, fmt.Errorf("Container with name %s not found", name)
}

// RoleContainer container for role steps. This inherites the context from
// the outer context. Any changes are lost once role has been completed
//
// Name - name of role
// File - file path that the role was loaded from
// User - user used to run all tasks with role
// Steps - list of steps to be ran against selected servers
type RoleContainer struct {
	Name  string
	File  string
	User  *hcl.Attribute
	Steps StepList
}
