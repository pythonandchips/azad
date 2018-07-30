package parser

import (
	multierror "github.com/hashicorp/go-multierror"
	"github.com/hashicorp/hcl2/hcl"
	"github.com/pythonandchips/azad/azad"
	"github.com/zclconf/go-cty/cty"
)

// PlaybookFromFile return a playbook
func PlaybookFromFile(path string, env map[string]string) (azad.Playbook, error) {
	errors := &multierror.Error{}
	evalContext := &hcl.EvalContext{
		Variables: map[string]cty.Value{},
	}
	if env, err := envToVariables(env); err == nil {
		evalContext.Variables["env"] = env
	}
	playbookDescription, err := loadConfigFile(path)
	if err != nil {
		errors = multierror.Append(errors, err)
	}
	variables, err := unpackVariables(playbookDescription.Variables, evalContext)
	if err != nil {
		return azad.Playbook{}, err
	}
	if len(variables) == 0 {
		evalContext.Variables["var"] = cty.MapValEmpty(cty.String)
	} else {
		evalContext.Variables["var"] = cty.MapVal(variables)
	}
	inventories, err := unpackInventories(playbookDescription.Inventories, evalContext)
	if err != nil {
		errors = multierror.Append(errors, err)
	}
	servers, err := unpackServer(playbookDescription.Servers)
	if err != nil {
		errors = multierror.Append(errors, err)
	}
	hosts, err := unpackHosts(playbookDescription.Hosts, playbookDescription.roleDescriptionGroups, evalContext)
	if err != nil {
		errors = multierror.Append(errors, err)
	}
	return azad.Playbook{
		Inventories: inventories,
		Servers:     servers,
		Hosts:       hosts,
		Variables:   variables,
	}, errors.ErrorOrNil()
}

func envToVariables(env map[string]string) (cty.Value, error) {
	variables := map[string]cty.Value{}
	for key, val := range env {
		variables[key] = cty.StringVal(val)
	}
	if len(variables) == 0 {
		return cty.MapValEmpty(cty.String), nil
	}
	return cty.MapVal(variables), nil
}
