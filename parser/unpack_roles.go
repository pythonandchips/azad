package parser

import (
	"fmt"

	multierror "github.com/hashicorp/go-multierror"
	"github.com/hashicorp/hcl2/hcl"
	"github.com/pythonandchips/azad/azad"
	"github.com/zclconf/go-cty/cty"
)

func unpackRoles(
	roleDescriptionGroups roleDescriptionGroups,
	evalContext *hcl.EvalContext,
) ([]azad.Role, error) {
	errors := &multierror.Error{}
	roles := []azad.Role{}
	for _, roleDescriptionGroup := range roleDescriptionGroups {
		role, err := unpackRole(roleDescriptionGroup, evalContext)
		if err != nil {
			errors = multierror.Append(errors, err)
		}
		roles = append(roles, role)
	}
	return roles, errors.ErrorOrNil()
}

func unpackRole(
	roleDescriptionGroup roleDescriptionGroup,
	evalContext *hcl.EvalContext,
) (azad.Role, error) {
	role := azad.Role{
		Name:      roleDescriptionGroup.name,
		Variables: map[string]cty.Value{},
		Path:      roleDescriptionGroup.path,
	}
	err := unpackTasksForRole(
		"main",
		roleDescriptionGroup,
		evalContext,
		&role,
	)
	if err != nil {
		return role, err
	}
	return role, nil
}

func unpackTasksForRole(
	file string,
	roleDescriptionGroup roleDescriptionGroup,
	evalContext *hcl.EvalContext,
	role *azad.Role,
) error {
	roleContext := newChildContext(evalContext)
	mainFile, err := roleDescriptionGroup.findFile(file)
	if err != nil {
		return fmt.Errorf("cannot find main file for role %s", roleDescriptionGroup.name)
	}
	roleDescription := mainFile.roleDescription
	user := roleDescription.User
	if role.User == "" {
		role.User = user
	}
	dependents := roleDescription.Dependents
	variables, err := unpackVariables(roleDescription.Variables, roleContext)
	if err != nil {
		return err
	}
	roleSchema := roleSchema()
	content, configErr := roleDescription.Config.Content(&roleSchema)
	if configErr.HasErrors() {
		return fmt.Errorf("unable to parse role config: %s", configErr.Error())
	}
	errors := &multierror.Error{}
	for _, block := range content.Blocks {
		switch block.Type {
		case "includes":
			err := unpackTaskBlock(
				block,
				roleContext,
				roleDescriptionGroup,
				role,
			)
			if err != nil {
				errors = multierror.Append(errors, err)
			}
		case "task":
			task, _ := unpackTask(block, roleContext)
			role.Tasks = append(role.Tasks, task)
		default:
			errors = multierror.Append(errors, fmt.Errorf("unrecognized block %s in %s/%s.az", block.Type, roleDescriptionGroup.name, file))
		}
	}
	role.Dependents = append(role.Dependents, dependents...)
	role.Variables = mergeMap(role.Variables, variables)
	return errors.ErrorOrNil()
}

func unpackTaskBlock(block *hcl.Block, roleContext *hcl.EvalContext, roleDescriptionGroup roleDescriptionGroup, role *azad.Role) error {
	switch block.Type {
	case "includes":
		includes, err := unpackIncludeMany(block, roleContext)
		if err != nil {
			return err
		}
		err = includeFilesInRole(
			includes,
			roleDescriptionGroup,
			roleContext,
			role,
		)
		if err != nil {
			return err
		}
	case "task":
		task, _ := unpackTask(block, roleContext)
		role.Tasks = append(role.Tasks, task)
	default:
		return fmt.Errorf("unrecognized block %s in %s", block.Type, roleDescriptionGroup.name)
	}
	return nil
}

func mergeMap(existing, new map[string]cty.Value) map[string]cty.Value {
	for key, value := range new {
		existing[key] = value
	}
	return existing
}

func includeFilesInRole(
	includes []string,
	roleDescriptionGroup roleDescriptionGroup,
	roleContext *hcl.EvalContext,
	role *azad.Role,
) error {
	errors := &multierror.Error{}
	for _, include := range includes {
		err := unpackTasksForRole(include, roleDescriptionGroup, roleContext, role)
		if err != nil {
			errors = multierror.Append(errors, err)
			continue
		}
	}
	return errors.ErrorOrNil()
}

func unpackIncludeMany(body *hcl.Block, evalContext *hcl.EvalContext) ([]string, error) {
	paths := []string{}
	attributesList, _ := body.Body.JustAttributes()
	pathAttribute, err := attributeFromList(attributesList, "paths")
	if err != nil {
		return paths, err
	}
	values, _ := pathAttribute.Expr.Value(evalContext)
	for _, value := range values.AsValueSlice() {
		paths = append(paths, value.AsString())
	}
	return paths, nil
}

func attributeFromList(attributeList hcl.Attributes, name string) (*hcl.Attribute, error) {
	for _, attr := range attributeList {
		if attr.Name == name {
			return attr, nil
		}
	}
	return nil, fmt.Errorf("Missing required attribute %s", name)

}

func unpackTask(body *hcl.Block, evalContext *hcl.EvalContext) (azad.Task, error) {
	task := azad.Task{
		Type: body.Labels[0],
		Name: body.Labels[1],
	}
	attributes := map[string]*hcl.Attribute{}
	attributesList, _ := body.Body.JustAttributes()
	for _, attr := range attributesList {
		if attr.Name == "condition" {
			task.Condition = attr
			continue
		}
		attributes[attr.Name] = attr
	}
	task.Attributes = attributes
	return task, nil
}

func roleSchema() hcl.BodySchema {
	roleSchema := hcl.BodySchema{
		Blocks: []hcl.BlockHeaderSchema{
			{Type: "includes"},
			{Type: "task", LabelNames: []string{"command", "label"}},
			{Type: "variable", LabelNames: []string{"name"}},
		},
	}
	return roleSchema
}
