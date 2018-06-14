package parser

import (
	"github.com/hashicorp/hcl2/hcl"
	"github.com/pythonandchips/azad/azad"
)

type hostDescription struct {
	ServerGroup string   `hcl:",label"`
	Config      hcl.Body `hcl:",remain"`
}

func unpackHosts(hostDescriptions []hostDescription, evalContext *hcl.EvalContext) ([]azad.Host, error) {
	hosts := []azad.Host{}
	hostSchema := hostSchema()
	for _, hostDescription := range hostDescriptions {
		hostSource, err := hostDescription.Config.Content(&hostSchema)
		var roles []string
		if attr, ok := hostSource.Attributes["roles"]; ok {
			val, _ := attr.Expr.Value(evalContext)
			for _, role := range val.AsValueSlice() {
				roles = append(roles, role.AsString())
			}
		}
		if err != nil {
			return hosts, err
		}
		host := azad.Host{
			ServerGroup: hostDescription.ServerGroup,
			Roles:       roles,
		}
		hosts = append(hosts, host)
	}
	return hosts, nil
}

func hostSchema() hcl.BodySchema {
	hostSchema := hcl.BodySchema{
		Attributes: []hcl.AttributeSchema{
			{Name: "roles", Required: true},
		},
	}
	return hostSchema
}
