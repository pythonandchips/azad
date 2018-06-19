package parser

import (
	"fmt"

	"github.com/hashicorp/hcl2/hcl"
)

type playbookDescription struct {
	Inventories []inventoryDescription `hcl:"inventory,block"`
	Servers     []serverDescription    `hcl:"server,block"`
	Hosts       []hostDescription      `hcl:"host,block"`
	Variables   []variableDescription  `hcl:"variable,block"`
	Roles       []roleDescription      `hcl:"role,block"`
}

type inventoryDescription struct {
	Name   string   `hcl:",label"`
	Config hcl.Body `hcl:",remain"`
}

type serverDescription struct {
	Group     string   `hcl:",label"`
	Addresses []string `hcl:"addresses"`
	Config    hcl.Body `hcl:",remain"`
}

// Variable value
type variableDescription struct {
	Name    string `hcl:",label"`
	Default string `hcl:"default"`
}

type roleDescriptions []roleDescription

func (roleDescriptions roleDescriptions) RoleFor(name string) (roleDescription, error) {
	for _, roleDescription := range roleDescriptions {
		if roleDescription.Name == name {
			return roleDescription, nil
		}
	}
	return roleDescription{}, fmt.Errorf("role with name %s not found", name)
}

func (roleDescriptions roleDescriptions) ReplaceWith(role roleDescription) {
	for i, roleDescription := range roleDescriptions {
		if roleDescription.Name == role.Name {
			roleDescriptions[i] = role
		}
	}
}

// Role list of task to be applied to host
type roleDescription struct {
	Name   string            `hcl:",label"`
	Tasks  []taskDescription `hcl:"task,block"`
	Config hcl.Body          `hcl:",remain"`
}

// Task run command via ssh
type taskDescription struct {
	Type   string   `hcl:"type,label"`
	Name   string   `hcl:"name,label"`
	Config hcl.Body `hcl:",remain"`
}
