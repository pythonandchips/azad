package parser

import "github.com/zclconf/go-cty/cty"

func unpackVariables(variableDescriptions []variableDescription) (map[string]cty.Value, error) {
	variables := map[string]cty.Value{}
	for _, variableDescription := range variableDescriptions {
		variables[variableDescription.Name] = cty.StringVal(variableDescription.Default)
	}
	return variables, nil
}
