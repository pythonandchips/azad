package runner

import (
	"testing"

	"github.com/hashicorp/hcl2/hcl"
	"github.com/pythonandchips/azad/expect"
	"github.com/pythonandchips/azad/plugin"
	"github.com/pythonandchips/azad/steps"
	"github.com/stretchr/testify/assert"
	"github.com/zclconf/go-cty/cty"
)

func TestHandleInput(t *testing.T) {
	getInput = func(pluginName, serviceName string) (plugin.Input, error) {
		return plugin.Input{
			Fields: []plugin.Field{
				{Name: "path", Required: true},
			},
			Run: func(plugin.InputContext) (map[string]string, error) {
				return map[string]string{
					"key1": "val1",
					"key2": "val2",
				}, nil
			},
		}, nil
	}
	store := store{
		variables: variables{},
		servers:   servers{},
	}
	inputStep := steps.InputStep{
		Type: "ini_file",
		Body: &TestBody{
			attributes: map[string]*hcl.Attribute{
				"path": testExpression("path", `"ini_file.ini"`),
			},
		},
	}
	err := handleInput(inputStep, &store)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	expect.EqualFatal(t, len(store.variables), 1)
	variable := store.variables[0].value.AsValueMap()
	assert.Equal(t, variable["key1"], cty.StringVal("val1"))
	assert.Equal(t, variable["key2"], cty.StringVal("val2"))

}
