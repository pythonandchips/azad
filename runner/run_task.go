package runner

import (
	"github.com/hashicorp/hcl2/hcl"
	"github.com/pythonandchips/azad/azad"
	"github.com/pythonandchips/azad/logger"
	"github.com/pythonandchips/azad/plugin"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
)

func runTask(task azad.Task, taskSchema plugin.Task, runner *runner) error {
	logger.Info("Applying %s:%s on %s", task.Type, task.Name, runner.Address)
	evalContext := &hcl.EvalContext{
		Variables: runner.toContext(),
		Functions: stdFunctions(),
	}
	if task.Condition != nil {
		result, err := task.Condition.Expr.Value(evalContext)
		if err != nil {
			logger.Error("Error evaluating conditional %s:%s on %s: %s", task.Type, task.Name, runner.Address, err)
			return err
		}
		if result.True() {
			logger.Info("Skipping %s:%s on %s, condition failed", task.Type, task.Name, runner.Address)
			return nil
		}
	}
	vars, err := varsForTask(taskSchema.Fields, task, evalContext)
	if err != nil {
		return err
	}
	context := plugin.NewContext(vars, runner.Conn)
	results, err := taskSchema.Run(context)
	if err != nil {
		logger.Error("Failed %s:%s on %s", task.Type, task.Name, runner.Address)
		logger.Error("Error: %s", err)
		return err
	}
	runner.setResult(task.Name, results)
	logger.Info("Success %s:%s on %s", task.Type, task.Name, runner.Address)
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
