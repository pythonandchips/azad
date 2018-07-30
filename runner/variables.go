package runner

import "github.com/zclconf/go-cty/cty"

func mergeVariables(vars ...map[string]cty.Value) map[string]cty.Value {
	variables := map[string]cty.Value{}
	for _, variableSet := range vars {
		for key, value := range variableSet {
			variables[key] = value
		}
	}
	return variables
}
