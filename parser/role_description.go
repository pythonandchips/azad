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

func (roleDescriptionGroups roleDescriptionGroups) FilterMany(names []string) (roleDescriptionGroups, error) {
	filteredRoleDescriptionGroups := []roleDescriptionGroup{}
	for _, name := range names {
		roleDescriptionGroup, err := roleDescriptionGroups.Filter(name)
		if err != nil {
			return roleDescriptionGroups, err
		}
		filteredRoleDescriptionGroups = append(filteredRoleDescriptionGroups, roleDescriptionGroup)
	}
	return filteredRoleDescriptionGroups, nil
}

func (roleDescriptionGroups roleDescriptionGroups) Filter(name string) (roleDescriptionGroup, error) {
	for _, roleDescriptionGroup := range roleDescriptionGroups {
		if roleDescriptionGroup.name == name {
			return roleDescriptionGroup, nil
		}
	}
	return roleDescriptionGroup{}, fmt.Errorf("role description group not found")
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
	Name       string                `hcl:",label"`
	Dependents []string              `hcl:"dependents,optional"`
	Variables  []variableDescription `hcl:"variable,block"`
	Config     hcl.Body              `hcl:",remain"`
}

// Task run command via ssh
type taskDescription struct {
	Type   string   `hcl:"type,label"`
	Name   string   `hcl:"name,label"`
	Config hcl.Body `hcl:",remain"`
}
