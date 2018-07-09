package runner

import (
	"fmt"
	"testing"

	"github.com/hashicorp/hcl2/hcl"
	"github.com/hashicorp/hcl2/hcl/hclsyntax"
	"github.com/pythonandchips/azad/azad"
	"github.com/pythonandchips/azad/logger"
	"github.com/pythonandchips/azad/plugin"
	"github.com/stretchr/testify/assert"
)

func TestRunTask(t *testing.T) {
	var suppliedContext plugin.Context
	runner := &runner{
		Address: "10.0.0.1",
		taskResults: map[string]taskResult{
			"previous_result": map[string]string{
				"ok": "ok",
			},
			"ruby-exists": map[string]string{
				"exists": "false",
			},
			"python-exists": map[string]string{
				"exists": "true",
			},
		},
		Variables: map[string]string{
			"dynamic_value": "dvalue",
		},
	}
	t.Run("with a successful task", func(t *testing.T) {
		testLogger := logger.TestLogger()
		pluginTask := plugin.Task{
			Fields: []plugin.Field{
				{Name: "static_value", Type: "string", Required: false},
				{Name: "dynamic_value", Type: "string", Required: false},
				{Name: "other_task_value", Type: "string", Required: false},
			},
			Run: func(context plugin.Context) (map[string]string, error) {
				suppliedContext = context
				return map[string]string{
					"return_value": "true",
				}, nil
			},
		}
		task := azad.Task{
			Type: "nil-task",
			Name: "nil-task",
			Attributes: map[string]*hcl.Attribute{
				"static_value":     testExpression("static_value", "value"),
				"dynamic_value":    testExpression("dynamic_value", "${ var.dynamic_value }"),
				"other_task_value": testExpression("other_task_value", "${ previous_result.ok }"),
			},
		}
		err := runTask(task, pluginTask, runner)
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		t.Run("calculates and passes the variables to the task", func(t *testing.T) {
			assert.Equal(t, suppliedContext.Get("static_value"), "value")
			assert.Equal(t, suppliedContext.Get("dynamic_value"), "dvalue")
			assert.Equal(t, suppliedContext.Get("other_task_value"), "ok")
		})
		t.Run("adds the new variables to the runners variables", func(t *testing.T) {
			_, ok := runner.taskResults["nil-task"]
			assert.Equal(t, ok, true, "expected varaibles to contain %s", "nil-task")
		})
		t.Run("writes output to log", func(t *testing.T) {
			assertLengthFatal(t, len(testLogger.Lines), 2)
			assert.Equal(t, testLogger.Lines[0], "INFO: Applying nil-task:nil-task on 10.0.0.1\n")
			assert.Equal(t, testLogger.Lines[1], "INFO: Success nil-task:nil-task on 10.0.0.1\n")
		})
	})
	t.Run("logs error if task fails", func(t *testing.T) {
		testLogger := logger.TestLogger()
		failingPluginTask := plugin.Task{
			Fields: []plugin.Field{},
			Run: func(context plugin.Context) (map[string]string, error) {
				return map[string]string{}, fmt.Errorf("some task error")
			},
		}
		failingTask := azad.Task{
			Type:       "failing-task",
			Name:       "failing-task",
			Attributes: map[string]*hcl.Attribute{},
		}
		err := runTask(failingTask, failingPluginTask, runner)
		if err == nil {
			t.Fatalf("expected an error but got none")
		}
		assert.Equal(t, testLogger.Lines[0], "INFO: Applying failing-task:failing-task on 10.0.0.1\n")
		assert.Equal(t, testLogger.Lines[1], "ERR: Failed failing-task:failing-task on 10.0.0.1\n")
		assert.Equal(t, testLogger.Lines[2], "ERR: Error: some task error\n")
	})
	t.Run("skips task if condition is meet", func(t *testing.T) {
		testLogger := logger.TestLogger()
		taskHasRan := false
		conditionalPluginTask := plugin.Task{
			Fields: []plugin.Field{},
			Run: func(context plugin.Context) (map[string]string, error) {
				taskHasRan = true
				return map[string]string{}, nil
			},
		}
		conditionalTask := azad.Task{
			Type:       "conditional-task",
			Name:       "conditional-task",
			Attributes: map[string]*hcl.Attribute{},
			Condition:  testExpression("condition", "${not(ruby-exists.exists)}"),
		}
		err := runTask(conditionalTask, conditionalPluginTask, runner)
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		assert.Equal(t, taskHasRan, false)
		assert.Equal(t, testLogger.Lines[0], "INFO: Applying conditional-task:conditional-task on 10.0.0.1\n")
		assert.Equal(t, testLogger.Lines[1], "INFO: Skipping conditional-task:conditional-task on 10.0.0.1, condition failed\n")
	})
}

func testExpression(name, value string) *hcl.Attribute {
	expr, parseErr := hclsyntax.ParseTemplate([]byte(value), "", hcl.Pos{Line: 1, Column: 1})
	if parseErr.HasErrors() {
		panic("unable to parse test string: " + parseErr.Error())
	}
	return &hcl.Attribute{
		Name: name, Expr: expr,
	}
}

func assertLengthFatal(t *testing.T, length, expected int) {
	if length != expected {
		t.Fatalf("expected %d log lines but got %d", expected, length)
	}
}
