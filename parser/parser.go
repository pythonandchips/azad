package parser

import (
	"errors"

	"github.com/hashicorp/hcl2/hcl"
	"github.com/pythonandchips/azad/azad"
	"github.com/zclconf/go-cty/cty"
)

var errNoEnv = errors.New("No environment variables")

// PlaybookFromFile return a playbook
func PlaybookFromFile(path string, env map[string]string) (azad.Playbook, error) {
	evalContext := &hcl.EvalContext{
		Variables: map[string]cty.Value{},
	}
	if env, err := envToVariables(env); err == nil {
		evalContext.Variables["env"] = env
	}
	playbookDescription, err := loadConfigFile(path)
	if err != nil {
		return azad.Playbook{}, err
	}
	variables, err := unpackVariables(playbookDescription.Variables)
	evalContext.Variables["var"] = cty.MapVal(variables)
	servers, err := unpackServer(playbookDescription.Servers)
	if err != nil {
		return azad.Playbook{}, err
	}
	hosts, err := unpackHosts(playbookDescription.Hosts, evalContext)
	roles, err := unpackRoles(playbookDescription.Roles, evalContext)
	return azad.Playbook{
		Servers: servers,
		Hosts:   hosts,
		Roles:   roles,
	}, nil
}

func envToVariables(env map[string]string) (cty.Value, error) {
	variables := map[string]cty.Value{}
	for key, val := range env {
		variables[key] = cty.StringVal(val)
	}
	if len(variables) == 0 {
		return cty.StringVal(""), errNoEnv
	}
	return cty.MapVal(variables), nil
}
