package runner

import (
	"testing"

	"github.com/hashicorp/hcl2/hcl"
	"github.com/pythonandchips/azad/conn"
	"github.com/pythonandchips/azad/logger"
	"github.com/pythonandchips/azad/plugin"
	"github.com/pythonandchips/azad/steps"
	"github.com/stretchr/testify/assert"
	"github.com/zclconf/go-cty/cty"
)

func TestHandleTask(t *testing.T) {
	logger.StubLogger()
	store := store{
		variables: variables{
			{name: "package", value: cty.StringVal("erlang-full"), valueType: "string"},
			{name: "has_installed", value: cty.StringVal("false"), valueType: "string"},
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
	t.Run("run the required task", func(t *testing.T) {
		getTask = func(pluginName, taskName string) (plugin.Task, error) {
			return testFullPluginTask(), nil
		}
		conn := &conn.LoggerSSHConn{}
		connection := connection{
			taskVariables: variables{
				{name: "install-erlang.success", value: cty.StringVal("true"), valueType: "string"},
			},
			conn: conn,
		}
		contextStore := contextStore{
			store:              store,
			contextUser:        "deploy",
			contextPath:        "/home/user/azad/roles/ruby",
			contextConnections: connections{&connection},
		}
		taskStep := steps.TaskStep{
			Type:  "apt-get.install",
			Label: "install-erlang",
			User:  testExpression("user", `"root"`),
			Body: &TestBody{
				attributes: map[string]*hcl.Attribute{
					"string":        testExpression("string", `"string"`),
					"variable":      testExpression("variable", `var.package`),
					"interpolation": testExpression("interpolation", `"${var.package}-name"`),
					"map-access":    testExpression("map-access", `var.folders["install"]`),
					"array-access":  testExpression("array-access", `var.packages[0]`),
					"composite":     testExpression("composite", `"${var.packages[0]} ${var.folders["install"]}, var.package"`),
				},
			},
		}
		err := handleTask(taskStep, contextStore)
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		assert.Equal(t, conn.Commands[0].User, "root")
		assert.Equal(t, conn.Commands[0].Command[0], "string")
		assert.Equal(t, conn.Commands[0].Command[1], "erlang-full")
		assert.Equal(t, conn.Commands[0].Command[2], "erlang-full-name")
		assert.Equal(t, conn.Commands[0].Command[3], "/opt")
		assert.Equal(t, conn.Commands[0].Command[4], "erlang-full")
		assert.Equal(t, connection.taskVariables[1].name, "install-erlang.installed")
	})
	t.Run("run the when condition passes", func(t *testing.T) {
		getTask = func(pluginName, taskName string) (plugin.Task, error) {
			return testFullPluginTask(), nil
		}
		conn := &conn.LoggerSSHConn{}
		connection := connection{
			taskVariables: variables{
				{name: "install-erlang.success", value: cty.StringVal("true"), valueType: "string"},
			},
			conn: conn,
		}
		contextStore := contextStore{
			store:              store,
			contextUser:        "deploy",
			contextPath:        "/home/user/azad/roles/ruby",
			contextConnections: connections{&connection},
		}
		taskStep := steps.TaskStep{
			Type:      "apt-get.install",
			Label:     "install-erlang",
			User:      testExpression("user", `"root"`),
			Condition: testExpression("condition", `var.has_installed == "true"`),
			Body: &TestBody{
				attributes: map[string]*hcl.Attribute{
					"string":        testExpression("string", `"string"`),
					"variable":      testExpression("variable", `var.package`),
					"interpolation": testExpression("interpolation", `"${var.package}-name"`),
					"map-access":    testExpression("map-access", `var.folders["install"]`),
					"array-access":  testExpression("array-access", `var.packages[0]`),
				},
			},
		}
		err := handleTask(taskStep, contextStore)
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		assert.Equal(t, len(conn.Commands), 1)
	})
	t.Run("run the when condition passes", func(t *testing.T) {
		getTask = func(pluginName, taskName string) (plugin.Task, error) {
			return testErrorPluginTask(), nil
		}
		conn := &conn.LoggerSSHConn{}
		connection := connection{
			taskVariables: variables{
				{name: "install-erlang.success", value: cty.StringVal("true"), valueType: "string"},
			},
			conn: conn,
		}
		contextStore := contextStore{
			store:              store,
			contextUser:        "deploy",
			contextPath:        "/home/user/azad/roles/ruby",
			contextConnections: connections{&connection},
		}
		taskStep := steps.TaskStep{
			Type:  "apt-get.install",
			Label: "install-erlang",
			User:  testExpression("user", `"root"`),
			Body: &TestBody{
				attributes: map[string]*hcl.Attribute{
					"string":        testExpression("string", `"string"`),
					"variable":      testExpression("variable", `var.package`),
					"interpolation": testExpression("interpolation", `"${var.package}-name"`),
					"map-access":    testExpression("map-access", `var.folders["install"]`),
					"array-access":  testExpression("array-access", `var.packages[0]`),
				},
			},
		}
		err := handleTask(taskStep, contextStore)
		if err == nil {
			t.Errorf("expected error but got none")
		}
	})
}
