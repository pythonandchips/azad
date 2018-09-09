package runner

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zclconf/go-cty/cty"
)

func TestContextStoreEvalContextFor(t *testing.T) {
	store := store{
		variables: variables{
			{name: "name", value: cty.StringVal("root"), valueType: "string"},
			{name: "aws.access_key", value: cty.StringVal("access_key"), valueType: "string"},
			{name: "aws.secret_key", value: cty.StringVal("secret_key"), valueType: "string"},
			{name: "folders", value: cty.MapVal(map[string]cty.Value{
				"install": cty.StringVal("/opt"),
			}), valueType: "map"},
			{name: "packages", value: cty.ListVal([]cty.Value{
				cty.StringVal("erlang-full"),
			}), valueType: "array"},
		},
		servers: servers{
			{address: "10.0.0.1", group: []string{"development", "kibana_server"}},
			{address: "10.0.0.2", group: []string{"kibana_server"}},
			{address: "10.0.0.3", group: []string{"development"}},
		},
	}
	connection := connection{
		taskVariables: variables{
			{name: "install-erlang.success", value: cty.StringVal("true"), valueType: "string"},
		},
	}
	contextStore := contextStore{store: store}
	assertExpressionForTask(t, contextStore, connection, "will_install", "install-erlang.success", "true")

}

func assertExpressionForTask(
	t *testing.T,
	store contextStore,
	connection connection,
	name,
	variable,
	value string,
) {
	expression := testExpression(name, variable)
	evalContext, err := store.evalContextForTask(expression.Expr.Variables(), &connection)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	evaluatedValue, valueErr := expression.Expr.Value(evalContext)
	if valueErr.HasErrors() {
		t.Fatalf("unexpected error: %s", valueErr)
	}
	assert.Equal(t, evaluatedValue.AsString(), value)
}
