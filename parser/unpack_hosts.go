package parser

import (
	"fmt"
	"strings"

	multierror "github.com/hashicorp/go-multierror"
	"github.com/hashicorp/hcl2/hcl"
	"github.com/pythonandchips/azad/azad"
	"github.com/pythonandchips/azad/helpers/stringslice"
	"github.com/zclconf/go-cty/cty"
)

type hostDescription struct {
	ServerGroup string                `hcl:",label"`
	Variables   []variableDescription `hcl:"variable,block"`
	Config      hcl.Body              `hcl:",remain"`
}

func unpackHosts(hostDescriptions []hostDescription, roleDescriptions roleDescriptionGroups, evalContext *hcl.EvalContext) ([]azad.Host, error) {
	hosts := []azad.Host{}
	hostSchema := hostSchema()
	errors := &multierror.Error{}
	for _, hostDescription := range hostDescriptions {
		hostContext := newChildContext(evalContext)
		hostSource, configErr := hostDescription.Config.Content(&hostSchema)
		if configErr != nil {
			errors = multierror.Append(errors, configErr)
			continue
		}
		variables, err := unpackVariables(hostDescription.Variables, hostContext)
		if err != nil {
			errors = multierror.Append(errors, err)
		}
		var roleNames []string
		if attr, ok := hostSource.Attributes["roles"]; ok {
			val, _ := attr.Expr.Value(evalContext)
			for _, role := range val.AsValueSlice() {
				roleNames = append(roleNames, role.AsString())
			}
		}
		roles, err := unpackRolesForHost(roleNames, roleDescriptions, hostContext, []string{})
		if err != nil {
			errors = multierror.Append(errors, err)
			continue
		}
		host := azad.Host{
			ServerGroup: hostDescription.ServerGroup,
			Variables:   variables,
			Roles:       roles,
		}
		hosts = append(hosts, host)
	}
	return hosts, errors.ErrorOrNil()
}

func unpackRolesForHost(
	roleNames []string,
	roleDescriptionGroups roleDescriptionGroups,
	evalContext *hcl.EvalContext,
	previousRoles []string,
) (azad.Roles, error) {
	roles := []azad.Role{}
	for _, roleName := range roleNames {
		if stringslice.Exists(roleName, previousRoles) {
			previousRoles = append(previousRoles, roleName)
			path := strings.Join(previousRoles, " > ")
			return roles, fmt.Errorf("dependent Loop detected %s", path)
		}
		previousRoles = append(previousRoles, roleName)
		roleDescriptionGroup, err := roleDescriptionGroups.RoleFor(roleName)
		if err != nil {
			return roles, err
		}
		role, err := unpackRole(roleDescriptionGroup, evalContext)
		if err != nil {
			return roles, err
		}
		dependentRoles, err := unpackRolesForHost(role.Dependents, roleDescriptionGroups, evalContext, previousRoles)
		if err != nil {
			return roles, err
		}
		roles = append(roles, dependentRoles...)
		roles = append(roles, role)
		previousRoles = previousRoles[:len(previousRoles)-1]
	}
	return roles, nil
}

func hostSchema() hcl.BodySchema {
	hostSchema := hcl.BodySchema{
		Attributes: []hcl.AttributeSchema{
			{Name: "roles", Required: true},
		},
	}
	return hostSchema
}

func newChildContext(evalContext *hcl.EvalContext) *hcl.EvalContext {
	childContext := evalContext.NewChild()
	childContext.Variables = map[string]cty.Value{
		"var": evalContext.Variables["var"],
	}
	return childContext
}
