package parser

import (
	multierror "github.com/hashicorp/go-multierror"
	"github.com/hashicorp/hcl2/hcl"
	"github.com/zclconf/go-cty/cty"
)

func unpackVariables(variableDescriptions []variableDescription, evalContext *hcl.EvalContext) (map[string]cty.Value, error) {
	variables := map[string]cty.Value{}
	variableSchema := variableSchema()
	errors := &multierror.Error{}
	for _, variableDescription := range variableDescriptions {
		variableSource, err := variableDescription.Config.Content(&variableSchema)
		if err != nil {
			errors = multierror.Append(errors, err)
			continue
		}
		if attr, ok := variableSource.Attributes["default"]; ok {
			val, err := attr.Expr.Value(evalContext)
			if err != nil {
				errors = multierror.Append(errors, err)
				continue
			}
			variables[variableDescription.Name] = val
			contextVars := evalContext.Variables["var"]
			vars := map[string]cty.Value{}
			if !contextVars.IsNull() {
				vars = contextVars.AsValueMap()
			}
			if len(vars) == 0 {
				vars = map[string]cty.Value{}
			}
			vars[variableDescription.Name] = val
			evalContext.Variables["var"] = cty.MapVal(vars)
		}
	}
	return variables, errors.ErrorOrNil()
}

func variableSchema() hcl.BodySchema {
	hostSchema := hcl.BodySchema{
		Attributes: []hcl.AttributeSchema{
			{Name: "default", Required: true},
		},
	}
	return hostSchema
}
