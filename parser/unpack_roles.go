package parser

import (
	"fmt"

	multierror "github.com/hashicorp/go-multierror"
	"github.com/hashicorp/hcl2/hcl"
	"github.com/pythonandchips/azad/azad"
)

func unpackRoles(roleDescriptionGroups roleDescriptionGroups, evalContext *hcl.EvalContext) ([]azad.Role, error) {
	errors := &multierror.Error{}
	roles := []azad.Role{}
	for _, roleDescriptionGroup := range roleDescriptionGroups {
		role := azad.Role{
			Name: roleDescriptionGroup.name,
		}
		tasks, dependents, err := unpackTasksForRole(
			"main",
			roleDescriptionGroup,
			evalContext,
		)
		if err != nil {
			errors = multierror.Append(errors, err)
			continue
		}
		role.Dependents = append(role.Dependents, dependents...)
		role.Tasks = tasks
		roles = append(roles, role)
	}
	return roles, errors.ErrorOrNil()
}

func unpackTasksForRole(
	file string,
	roleDescriptionGroup roleDescriptionGroup,
	evalContext *hcl.EvalContext,
) ([]azad.Task, []string, error) {
	tasks := []azad.Task{}
	mainFile, err := roleDescriptionGroup.findFile(file)
	if err != nil {
		return tasks, []string{}, fmt.Errorf("cannot find main file for role %s", roleDescriptionGroup.name)
	}
	roleDescription := mainFile.roleDescription
	dependents := roleDescription.Dependents
	roleSchema := roleSchema()
	content, configErr := roleDescription.Config.Content(&roleSchema)
	if configErr.HasErrors() {
		return tasks, []string{}, fmt.Errorf("unable to parse role config: %s", err.Error())
	}
	errors := &multierror.Error{}
	for _, block := range content.Blocks {
		switch block.Type {
		case "include":
			include := unpackInclude(block, evalContext)
			includeTasks, includeDependents, err := unpackTasksForRole(include, roleDescriptionGroup, evalContext)
			if err != nil {
				errors = multierror.Append(errors, err)
				continue
			}
			tasks = append(tasks, includeTasks...)
			dependents = append(dependents, includeDependents...)
		case "task":
			task, _ := unpackTask(block, evalContext)
			tasks = append(tasks, task)
		default:
			errors = multierror.Append(errors, fmt.Errorf("unrecognized block %s in %s/%s.az", block.Type, roleDescriptionGroup.name, file))
		}
	}
	return tasks, dependents, nil
}

func unpackInclude(body *hcl.Block, evalContext *hcl.EvalContext) string {
	attributes := map[string]string{}
	attributesList, _ := body.Body.JustAttributes()
	for _, attr := range attributesList {
		value, _ := attr.Expr.Value(evalContext)
		attributes[attr.Name] = value.AsString()
	}
	return attributes["path"]
}

func unpackTask(body *hcl.Block, evalContext *hcl.EvalContext) (azad.Task, error) {
	task := azad.Task{
		Type: body.Labels[0],
		Name: body.Labels[1],
	}
	attributes := map[string]*hcl.Attribute{}
	attributesList, _ := body.Body.JustAttributes()
	for _, attr := range attributesList {
		attributes[attr.Name] = attr
	}
	task.Attributes = attributes
	return task, nil
}

func roleSchema() hcl.BodySchema {
	roleSchema := hcl.BodySchema{
		Blocks: []hcl.BlockHeaderSchema{
			{Type: "include"},
			{Type: "task", LabelNames: []string{"command", "label"}},
		},
	}
	return roleSchema
}
