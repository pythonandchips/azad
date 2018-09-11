package runner

import (
	"github.com/pythonandchips/azad/communicator"
	"github.com/pythonandchips/azad/plugin"
	"github.com/pythonandchips/azad/steps"
	"github.com/zclconf/go-cty/cty"
)

func handleInput(inputStep steps.InputStep, inputStore inputStore) error {
	inputSchema, err := inputSchema(inputStep)
	if err != nil {
		return err
	}
	attributes, err := attributesForSchema(inputStep.Body, inputSchema.Fields, inputStore)
	if err != nil {
		return err
	}
	inputContext := plugin.NewInputContext(attributes, inputStore.playbookPath(), inputStore.rolePath())
	variables, err := inputSchema.Run(inputContext)

	variableMap := map[string]cty.Value{}
	for key, value := range variables {
		variableMap[key] = cty.StringVal(value)
	}
	inputStore.addVariable(variable{
		name:      inputStep.Name,
		value:     variableMapFromMap(variableMap),
		valueType: "map",
	})
	return nil
}

func variableMapFromMap(variableMap map[string]cty.Value) cty.Value {
	if len(variableMap) == 0 {
		return cty.MapValEmpty(cty.String)
	}
	return cty.MapVal(variableMap)
}

func inputSchema(inputStep steps.InputStep) (plugin.Input, error) {
	pluginName := inputStep.PluginName()
	serviceName := inputStep.ServiceName()
	return getInput(pluginName, serviceName)
}

var getInput = func(pluginName, serviceName string) (plugin.Input, error) {
	return communicator.GetInput(pluginName, serviceName)
}
