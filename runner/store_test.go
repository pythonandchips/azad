package runner

import (
	"testing"

	"github.com/pythonandchips/azad/conn"
	"github.com/pythonandchips/azad/expect"
	"github.com/pythonandchips/azad/logger"
	"github.com/stretchr/testify/assert"
	"github.com/zclconf/go-cty/cty"
)

func TestStoreAddServer(t *testing.T) {
	t.Run("add new server", func(t *testing.T) {
		servers := servers{
			{address: "10.0.0.1", group: []string{"kibana_server"}},
		}

		store := store{}
		store.addServers(servers)

		expect.EqualFatal(t, len(store.servers), 1)
		server := store.servers[0]

		assert.Equal(t, server.address, "10.0.0.1")
		assert.Equal(t, server.group, []string{"kibana_server"})
	})

	t.Run("add server that already exists with different group", func(t *testing.T) {
		store := store{
			servers: servers{
				{address: "10.0.0.1", group: []string{"kibana_server"}},
			},
		}

		servers := servers{
			{address: "10.0.0.1", group: []string{"group_development"}},
		}

		store.addServers(servers)

		expect.EqualFatal(t, len(store.servers), 1)
		server := store.servers[0]

		assert.Equal(t, server.address, "10.0.0.1")
		assert.Equal(t, server.group, []string{"kibana_server", "group_development"})
	})

	t.Run("adds server that already exists with the given group", func(t *testing.T) {
		store := store{
			servers: servers{
				{address: "10.0.0.1", group: []string{"kibana_server"}},
			},
		}
		servers := servers{
			{address: "10.0.0.1", group: []string{"kibana_server"}},
		}

		store.addServers(servers)

		expect.EqualFatal(t, len(store.servers), 1)
		server := store.servers[0]

		assert.Equal(t, server.address, "10.0.0.1")
		assert.Equal(t, server.group, []string{"kibana_server"})
	})
}

func TestContextStore(t *testing.T) {
	logger.StubLogger()
	conn.Simulate()
	store := store{
		variables: variables{
			{name: "log-path", value: cty.StringVal("/var/log"), valueType: "string"},
		},
		servers: servers{
			{address: "10.0.0.1", group: []string{"kibana_server"}},
		},
	}
	contextStore, err := store.contextStore(
		[]string{"kibana_server"},
		"admin",
		"roles/kibana_server",
	)
	if err != nil {
		t.Fatalf("unexpected err: %s", err)
	}
	expect.EqualFatal(t, len(contextStore.variables), 1)
	assert.Equal(t, contextStore.variables[0].name, "log-path")
	assert.Equal(t, contextStore.variables[0].value, cty.StringVal("/var/log"))
	assert.Equal(t, contextStore.variables[0].valueType, "string")
	assert.Equal(t, len(contextStore.contextConnections), 1)
	assert.Equal(t, contextStore.contextUser, "admin")
	assert.Equal(t, contextStore.contextPath, "roles/kibana_server")
	t.Run("updating context store does not change main store", func(t *testing.T) {
		contextStore.addVariable(variable{name: "other-log-path"})

		assert.Equal(t, len(contextStore.variables), 2)
		assert.Equal(t, len(store.variables), 1)
	})
}

func TestStoreContextFor(t *testing.T) {
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

	t.Run("access string variable", func(t *testing.T) {
		assertExpression(t, store, "body", `var.name`, "root")
	})
	t.Run("access map variable", func(t *testing.T) {
		assertExpression(t, store, "install_folder", `var.folders["install"]`, "/opt")
	})
	t.Run("access array value", func(t *testing.T) {
		assertExpression(t, store, "install_folder", `var.packages[0]`, "erlang-full")
	})
	t.Run("access compound variable", func(t *testing.T) {
		assertExpression(t, store, "install_folder", `var.aws.access_key`, "access_key")
	})
	t.Run("access compound variable with multiple keys", func(t *testing.T) {
		assertExpression(t, store, "install_folder", `"${var.aws.access_key}-${var.aws.secret_key}"`, "access_key-secret_key")
	})
	t.Run("access server addresses", func(t *testing.T) {
		assertExpression(t, store, "db_server", `srv.kibana_server[0]`, "10.0.0.1")
		assertExpression(t, store, "db_server", `srv.kibana_server[1]`, "10.0.0.2")
	})
	t.Run("can resolve boolean", func(t *testing.T) {
		expression := testExpression("condition", `var.name == "root"`)
		evalContext, err := store.evalContextFor(expression.Expr.Variables())
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		evaluatedValue, valueErr := expression.Expr.Value(evalContext)
		if valueErr.HasErrors() {
			t.Fatalf("unexpected error: %s", valueErr)
		}
		assert.Equal(t, evaluatedValue.True(), true)
	})
}

func TestChildStore(t *testing.T) {
	store := store{
		variables: variables{
			{name: "log-path", value: cty.StringVal("/var/log"), valueType: "string"},
		},
		servers: servers{
			{address: "10.0.0.1", group: []string{"kibana_server"}},
		},
	}
	contextStore := contextStore{
		store:       store,
		contextUser: "admin",
		contextConnections: connections{
			{
				taskVariables: variables{
					{name: "log-folder"},
				},
			},
		},
	}
	newContextStore := contextStore.childStore("root", "roles/new_context")
	assert.Equal(t, newContextStore.contextUser, "root")
	assert.Equal(t, newContextStore.contextPath, "roles/new_context")
	t.Run("does not change parent context", func(t *testing.T) {
		newContextStore.addVariable(variable{name: "log-service"})

		assert.Equal(t, len(newContextStore.variables), 2)
		assert.Equal(t, len(contextStore.variables), 1)
	})
}

func assertExpression(t *testing.T, store store, name, variable, value string) {
	expression := testExpression(name, variable)
	evalContext, err := store.evalContextFor(expression.Expr.Variables())
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	evaluatedValue, valueErr := expression.Expr.Value(evalContext)
	if valueErr.HasErrors() {
		t.Fatalf("unexpected error: %s", valueErr)
	}
	assert.Equal(t, evaluatedValue.AsString(), value)
}
