package runner

import (
	multierror "github.com/hashicorp/go-multierror"
	"github.com/pythonandchips/azad/azad"
	"github.com/pythonandchips/azad/communicator"
	"github.com/pythonandchips/azad/conn"
	"github.com/pythonandchips/azad/logger"
	"github.com/pythonandchips/azad/parser"
)

var config azad.Config

// RunPlaybook run the playbook
var RunPlaybook = func(playbookFilePath string, globalConfig azad.Config) {
	config = globalConfig
	env := map[string]string{}
	playbook, err := parser.PlaybookFromFile(playbookFilePath, env)
	if err != nil {
		logger.ErrorAndExit("Playbook is invalid: %s", err)
	}
	playbook, err = readInventory(playbook)
	if err != nil {
		logger.ErrorAndExit("unable to load inventory: %s", err)
		return
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
		connection := conn.NewConn(config.SSHConfig())
		err := connection.ConnectTo(address)
		if err != nil {
			errors = multierror.Append(errors, err)
			break
		}
		runner := newRunner(connection, address)
		runners = append(runners, runner)
	}
	return runners, errors.ErrorOrNil()
}

// ValidatePlugins validate that all plugins and tasks are available
func validatePlugins(playbook azad.Playbook) error {
	errors := &multierror.Error{}
	tasks := playbook.RequiredTasks()
	for _, task := range tasks {
		if _, err := communicator.GetTask(task.PluginName(), task.TaskName()); err != nil {
			errors = multierror.Append(errors, err)
		}
	}
	return errors.ErrorOrNil()
}

func runTasks(tasks azad.Tasks, runners runners) error {
	for _, task := range tasks {
		taskSchema, _ := communicator.GetTask(task.PluginName(), task.TaskName())
		for _, runner := range runners {
			runTask(task, taskSchema, runner)
		}
	}
	return nil
}
