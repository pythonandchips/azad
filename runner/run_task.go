package runner

import (
	"errors"

	"github.com/hashicorp/hcl2/hcl"
	"github.com/pythonandchips/azad/azad"
	"github.com/pythonandchips/azad/conn"
	"github.com/pythonandchips/azad/logger"
	"github.com/pythonandchips/azad/plugin"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
)

type runTaskParams struct {
	task       azad.Task
	taskSchema plugin.Task
	runner     *runner
	rootPath   string
	rolePath   string
}

func (runTaskParams runTaskParams) Type() string {
	return runTaskParams.task.Type
}

func (runTaskParams runTaskParams) Name() string {
	return runTaskParams.task.Name
}

func (runTaskParams runTaskParams) Address() string {
	return runTaskParams.runner.Address
}

func (runTaskParams runTaskParams) toContext() map[string]cty.Value {
	return runTaskParams.runner.toContext()
}

func (runTaskParams runTaskParams) conditionResult(evalContext *hcl.EvalContext) (cty.Value, error) {
	value, err := runTaskParams.task.Condition.Expr.Value(evalContext)
	if err.HasErrors() {
		return value, errors.New(err.Error())
	}
	return value, nil
}

func (runTaskParams runTaskParams) conn() conn.Conn {
	return runTaskParams.runner.Conn
}

func (runTaskParams runTaskParams) run(context plugin.Context) (map[string]string, error) {
	return runTaskParams.taskSchema.Run(context)
}

func (runTaskParams runTaskParams) logParmas() []interface{} {
	return []interface{}{
		runTaskParams.Type(),
		runTaskParams.Name(),
		runTaskParams.Address(),
	}
}

func runTask(runTaskParams runTaskParams) error {
	logger.Info("Applying %s:%s on %s", runTaskParams.logParmas()...)
	evalContext := &hcl.EvalContext{
		Variables: runTaskParams.toContext(),
		Functions: stdFunctions(),
	}
	if runTaskParams.task.Condition != nil {
		result, err := runTaskParams.conditionResult(evalContext)
		if err != nil {
			logger.Error("Error evaluating conditional %s:%s on %s: %s", runTaskParams.logParmas()...)
			return err
		}
		if result.True() {
			logger.Info("Skipping %s:%s on %s, condition failed", runTaskParams.logParmas()...)
			return nil
		}
	}
	vars, err := varsForTask(runTaskParams.taskSchema.Fields, runTaskParams.task, evalContext)
	if err != nil {
		return err
	}
	context := plugin.NewContext(
		vars,
		runTaskParams.conn(),
		runTaskParams.rootPath,
		runTaskParams.rolePath,
	)
	results, err := runTaskParams.run(context)
	if err != nil {
		logger.Error("Failed %s:%s on %s", runTaskParams.logParmas()...)
		logger.Error("Error: %s", err)
		return err
	}
	runTaskParams.runner.setResult(runTaskParams.Name(), results)
	logger.Info("Success %s:%s on %s", runTaskParams.logParmas()...)
	return nil
}

func varsForTask(fields []plugin.Field, task azad.Task, evalContext *hcl.EvalContext) (map[string]string, error) {
	vars := map[string]string{}
	for _, field := range fields {
		attr := task.Attributes[field.Name]
		value, err := attr.Expr.Value(evalContext)
		if err != nil {
			return vars, err
		}
		vars[field.Name] = value.AsString()
	}
	return vars, nil
}

func stdFunctions() map[string]function.Function {
	return map[string]function.Function{
		"not": function.New(&function.Spec{
			Params: []function.Parameter{
				{
					Name:             "val",
					Type:             cty.String,
					AllowDynamicType: true,
				},
			},
			Type: function.StaticReturnType(cty.Bool),
			Impl: func(args []cty.Value, retType cty.Type) (cty.Value, error) {
				val := args[0] != cty.StringVal("true")
				return cty.BoolVal(val), nil
			},
		}),
		"is": function.New(&function.Spec{
			Params: []function.Parameter{
				{
					Name:             "val",
					Type:             cty.String,
					AllowDynamicType: true,
				},
			},
			Type: function.StaticReturnType(cty.Bool),
			Impl: func(args []cty.Value, retType cty.Type) (cty.Value, error) {
				val := args[0] == cty.StringVal("true")
				return cty.BoolVal(val), nil
			},
		}),
	}
}
