package runner

import (
	"github.com/pythonandchips/azad/steps"
	"github.com/zclconf/go-cty/cty"
)

func handleVariable(variableStep steps.VariableStep, store variableStore) error {
	variableType := cty.StringVal("string")
	if variableStep.Type != nil {
		var evalErr error
		variableType, evalErr = store.evalVariable(variableStep.Type, allowString)
		if evalErr != nil {
			return evalErr
		}
	}
	variableValue, evalErr := store.evalVariable(variableStep.Value, allowAll)
	if evalErr != nil {
		return evalErr
	}
	store.addVariable(variable{
		name:      variableStep.Name,
		value:     variableValue,
		valueType: variableType.AsString(),
	})
	return nil
}
