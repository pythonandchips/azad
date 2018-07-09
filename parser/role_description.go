package parser

import (
	"fmt"

	"github.com/hashicorp/hcl2/hcl"
)

type roleDescriptionGroups []roleDescriptionGroup

func (roleDescriptionGroups roleDescriptionGroups) RoleFor(name string) (roleDescriptionGroup, error) {
	for _, roleDescriptionGroup := range roleDescriptionGroups {
		if roleDescriptionGroup.name == name {
			return roleDescriptionGroup, nil
		}
	}
	return roleDescriptionGroup{}, fmt.Errorf("role with name %s not found", name)
}

func (roleDescriptionGroups roleDescriptionGroups) ReplaceWith(role roleDescriptionGroup) {
	for i, roleDescriptionGroup := range roleDescriptionGroups {
		if roleDescriptionGroup.name == role.name {
			roleDescriptionGroups[i] = role
		}
	}
}

type roleDescriptionGroup struct {
	name  string
	files []roleFile
}

func (roleDescriptionGroup roleDescriptionGroup) findFile(name string) (roleFile, error) {
	for _, roleFile := range roleDescriptionGroup.files {
		if roleFile.name == name {
			return roleFile, nil
		}
	}
	return roleFile{}, fmt.Errorf("role file with name %s not found", name)
}

type roleFile struct {
	name            string
	roleDescription roleDescription
}

type roleDescriptions []roleDescription

// Role list of task to be applied to host
type roleDescription struct {
	Name       string   `hcl:",label"`
	Dependents []string `hcl:"dependents,optional"`
	Config     hcl.Body `hcl:",remain"`
}

// Task run command via ssh
type taskDescription struct {
	Type   string   `hcl:"type,label"`
	Name   string   `hcl:"name,label"`
	Config hcl.Body `hcl:",remain"`
}
