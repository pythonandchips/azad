package parser

import (
	"github.com/hashicorp/hcl2/hcl"
	"github.com/pythonandchips/azad/azad"
)

func unpackRoles(rolesDescriptions []roleDescription, evalContext *hcl.EvalContext) ([]azad.Role, error) {
	roles := []azad.Role{}
	for _, roleDescription := range rolesDescriptions {
		role := azad.Role{
			Name:       roleDescription.Name,
			Dependents: roleDescription.Dependents,
		}
		tasks, _ := unpackTasks(roleDescription.Tasks, evalContext)
		role.Tasks = tasks
		roles = append(roles, role)
	}
	return roles, nil
}

func unpackTasks(taskDescriptions []taskDescription, evalContext *hcl.EvalContext) ([]azad.Task, error) {
	tasks := []azad.Task{}
	for _, taskDescription := range taskDescriptions {
		task := azad.Task{
			Type: taskDescription.Type,
			Name: taskDescription.Name,
		}
		attributes := map[string]string{}
		attributesList, _ := taskDescription.Config.JustAttributes()

		for _, attr := range attributesList {
			value, _ := attr.Expr.Value(evalContext)
			attributes[attr.Name] = value.AsString()
		}
		task.Attributes = attributes
		tasks = append(tasks, task)
	}
	return tasks, nil
}
