package runner

import (
	multierror "github.com/hashicorp/go-multierror"
	"github.com/pythonandchips/azad/azad"
	"github.com/pythonandchips/azad/conn"
	"github.com/pythonandchips/azad/logger"
	"github.com/pythonandchips/azad/parser"
	"github.com/pythonandchips/azad/plugins"
	"github.com/pythonandchips/azad/schema"
)

type runners []runner

func (runners runners) Close() {
	for _, runner := range runners {
		runner.Conn.Close()
	}
}

type runner struct {
	Address   string
	Conn      conn.Conn
	Variables map[string]string
}

// RunPlaybook run the playbook
func RunPlaybook(playbookFilePath string, config azad.Config) {
	env := map[string]string{}
	playbook, err := parser.PlaybookFromFile(playbookFilePath, env)
	if err != nil {
		logger.ErrorAndExit("Playbook is invalid: %s", err)
	}
	err = validatePlugins(playbook)
	if err != nil {
		logger.ErrorAndExit("Unable to load plugins: %s", err)
		return
	}
	for _, host := range playbook.Hosts {
		logger.Info("Running for %s", host.ServerGroup)
		tasks, err := playbook.TasksForRoles(host.Roles)

		if err != nil {
			logger.ErrorAndExit("%s", err)
		}

		server, _ := playbook.LookupServer(host.ServerGroup)
		runners, err := createRunners(server.Addresses)
		defer runners.Close()
		if err != nil {
			logger.Error("%s", err)
			return
		}
		runTasks(tasks, runners)
		logger.Info("Finished Applying %s", host.ServerGroup)
	}
}

func createRunners(addresses []string) (runners, error) {
	runners := runners{}
	errors := &multierror.Error{}
	for _, address := range addresses {
		connection := conn.NewConn()
		err := connection.ConnectTo(address)
		if err != nil {
			errors = multierror.Append(errors, err)
			break
		}
		runner := runner{Conn: connection, Address: address}
		runners = append(runners, runner)
	}
	return runners, errors.ErrorOrNil()
}

// ValidatePlugins validate that all plugins and tasks are available
func validatePlugins(playbook azad.Playbook) error {
	errors := &multierror.Error{}
	pluginLoader := plugins.Loader()
	tasks := playbook.RequiredTasks()
	for _, task := range tasks {
		if err := pluginLoader.TaskExists(task.PluginName(), task.TaskName()); err != nil {
			errors = multierror.Append(errors, err)
		}
	}
	return errors.ErrorOrNil()
}

func runTasks(tasks azad.Tasks, runners runners) error {
	for _, task := range tasks {
		taskSchema, _ := plugins.Loader().GetTask(task.PluginName(), task.TaskName())
		for _, runner := range runners {
			runTask(task, taskSchema, runner)
		}
	}
	return nil
}

func runTask(task azad.Task, taskSchema schema.Task, runner runner) error {
	logger.Info("Applying %s:%s on %s", task.Type, task.Name, runner.Address)
	context := schema.NewContext(task.Attributes, runner.Conn)
	err := taskSchema.Run(context)
	if err != nil {
		logger.Error("Failed %s:%s on %s", task.Type, task.Name, runner.Address)
		logger.Error("Error: %s", err)
		return err
	}
	logger.Info("Success %s:%s on %s", task.Type, task.Name, runner.Address)
	return nil
}
