package runner

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zclconf/go-cty/cty"
)

func TestRunnerNewChild(t *testing.T) {
	runner := &runner{
		Address: "10.0.0.1",
		Conn:    &fakeConn{Closed: false},
		GlobalVariables: map[string]cty.Value{
			"install_dir": cty.StringVal("/opt"),
		},
	}
	childVariables := map[string]cty.Value{
		"install_ruby": cty.StringVal("/opt/ruby"),
	}
	childRunner := runner.newChild(childVariables)

	assert.Equal(t, childRunner.GlobalVariables["install_dir"].AsString(), "/opt")
	assert.Equal(t, childRunner.GlobalVariables["install_ruby"].AsString(), "/opt/ruby")
	assert.Equal(t, childRunner.Address, runner.Address)
	assert.Equal(t, childRunner.Conn, runner.Conn)
}

func TestRunnerToContext(t *testing.T) {
	runner := &runner{
		Address: "10.0.0.1",
		Conn:    &fakeConn{Closed: false},
		GlobalVariables: map[string]cty.Value{
			"install_dir": cty.StringVal("/opt"),
		},
		taskResults: map[string]taskResult{
			"ruby.install": taskResult{
				"success": "true",
			},
		},
	}
	context := runner.toContext()

	varValue := context["var"].AsValueMap()
	assert.Equal(t, varValue["install_dir"].AsString(), "/opt")

	rubyInstallResult := context["ruby.install"].AsValueMap()
	assert.Equal(t, rubyInstallResult["success"].AsString(), "true")

}
