package runner

import (
	"github.com/hashicorp/hcl2/hcl"
	"github.com/pythonandchips/azad/communicator"
	"github.com/pythonandchips/azad/logger"
	"github.com/pythonandchips/azad/plugin"
	"github.com/pythonandchips/azad/steps"
	"github.com/zclconf/go-cty/cty"
)

type runContext struct {
	taskStep   steps.TaskStep
	taskStore  taskStore
	connection *connection
	content    *hcl.BodyContent
	user       string
	taskSchema plugin.Task
}

var getTask = func(pluginName, taskName string) (plugin.Task, error) {
	return communicator.GetTask(pluginName, taskName)
}

func handleTask(taskStep steps.TaskStep, taskStore taskStore) error {
	taskSchema, err := getTask(taskStep.PluginName(), taskStep.TaskName())
	if err != nil {
		return err
	}
	taskBodySchema := createSchemaFromFields(taskSchema.Fields)
	content, contentErr := taskStep.Body.Content(taskBodySchema)
	if contentErr.HasErrors() {
		return contentErr
	}
	user, err := taskUser(taskStore, taskStep)
	if err != nil {
		return err
	}
	return taskStore.connections().each(func(connection *connection) error {
		runContext := runContext{
			taskStep:   taskStep,
			taskStore:  taskStore,
			connection: connection,
			content:    content,
			user:       user,
			taskSchema: taskSchema,
		}
		return runTaskOnConnection(runContext)
	})
}

func taskUser(taskStore taskStore, taskStep steps.TaskStep) (string, error) {
	if taskStep.User != nil {
		userVal, evalErr := taskStore.evalVariable(taskStep.User, allowString)
		if evalErr != nil {
			return "", evalErr
		}
		if userVal.AsString() != "" {
			return userVal.AsString(), nil
		}
	}
	return taskStore.user(), nil
}

func runTaskOnConnection(runContext runContext) error {
	logger.Info("Applying %s:%s on %s", runContext.taskStep.PluginName(), runContext.taskStep.TaskName(), runContext.connection.conn.Address())
	pass, err := runContext.shouldRun()
	if err != nil {
		return err
	}
	if pass {
		logger.Info("Skipping %s:%s on %s, condition failed", runContext.taskStep.PluginName(), runContext.taskStep.TaskName(), runContext.connection.conn.Address())
		return nil
	}
	attrs, evalErr := runContext.taskAttributes()
	if evalErr != nil {
		return evalErr
	}

	context := plugin.NewContext(
		attrs,
		runContext.connection.conn,
		runContext.user,
		runContext.taskStore.playbookPath(),
		runContext.taskStore.rolePath(),
	)
	if runContext.debug() {
		logger.Debug("Input for %s:%s on %s", runContext.taskStep.PluginName(), runContext.taskStep.TaskName(), runContext.connection.conn.Address())
		logger.Debug("Attrs:")
		for k, v := range attrs {
			logger.Debug("  %s: %s", k, v.GoString())
		}
		logger.Debug("User: %s", runContext.user)
		logger.Debug("Playbook Path: %s", runContext.taskStore.playbookPath())
		logger.Debug("Role Path: %s", runContext.taskStore.rolePath())
	}
	results, err := runContext.taskSchema.Run(context)
	if runContext.debug() {
		logger.Debug("Result for %s:%s on %s", runContext.taskStep.PluginName(), runContext.taskStep.TaskName(), runContext.connection.conn.Address())
		logger.Debug("Attrs:")
		for k, v := range results {
			logger.Debug("  %s: %s", k, v)
		}
	}
	if err != nil {
		logger.Error("Failed %s:%s on %s, condition failed", runContext.taskStep.PluginName(), runContext.taskStep.TaskName(), runContext.connection.conn.Address())
		return err
	}
	for key, value := range results {
		runContext.connection.addTaskVariable(variable{
			name:      runContext.taskStep.Label + "." + key,
			valueType: "string",
			value:     cty.StringVal(value),
		})
	}
	logger.Info("Success %s:%s on %s", runContext.taskStep.PluginName(), runContext.taskStep.TaskName(), runContext.connection.conn.Address())
	return nil
}

func (runContext runContext) shouldRun() (bool, error) {
	if runContext.taskStep.Condition != nil {
		val, evalErr := runContext.taskStore.evalVariableForTask(runContext.taskStep.Condition, runContext.connection, allowAll)
		if evalErr != nil {
			return false, evalErr
		}
		return val.True(), nil
	}
	return false, nil
}

func (runContext runContext) debug() bool {
	if runContext.taskStep.Debug != nil {
		val, evalErr := runContext.taskStore.evalVariableForTask(runContext.taskStep.Debug, runContext.connection, allowAll)
		if evalErr != nil {
			return true
		}
		return val.True()
	}
	return false
}

func (runContext runContext) taskAttributes() (map[string]cty.Value, error) {
	attrs := map[string]cty.Value{}
	for name, attr := range runContext.content.Attributes {
		val, evalErr := runContext.taskStore.evalVariableForTask(attr, runContext.connection, allowAll)
		if evalErr != nil {
			return attrs, evalErr
		}
		attrs[name] = val
	}
	return attrs, nil
}
