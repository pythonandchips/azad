package runner

import (
	"github.com/hashicorp/hcl2/hcl"
	"github.com/pythonandchips/azad/azad"
	"github.com/pythonandchips/azad/logger"
	"github.com/pythonandchips/azad/plugin"
)

func runTask(task azad.Task, taskSchema plugin.Task, runner *runner) error {
	logger.Info("Applying %s:%s on %s", task.Type, task.Name, runner.Address)
	vars := map[string]string{}
	for _, field := range taskSchema.Fields {
		evalContext := &hcl.EvalContext{
			Variables: runner.toContext(),
		}
		attr := task.Attributes[field.Name]
		value, err := attr.Expr.Value(evalContext)
		if err != nil {
			return err
		}
		vars[field.Name] = value.AsString()
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
