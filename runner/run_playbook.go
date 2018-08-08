package runner

import (
	"time"

	"github.com/pythonandchips/azad/azad"
	"github.com/pythonandchips/azad/communicator"
	"github.com/pythonandchips/azad/logger"
	"github.com/pythonandchips/azad/parser"
)

var config azad.Config

var now = func() time.Time {
	return time.Now()
}

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

		if err != nil {
			logger.ErrorAndExit("%s", err)
		}

		server, _ := playbook.LookupServer(host.ServerGroup)
		runners, err := createRunners(server.Addresses, mergeVariables(playbook.Variables, host.Variables))
		defer runners.Close()
		if err != nil {
			logger.Error("%s", err)
			return
		}
		for _, role := range host.Roles {
			runTasks(role, runners.ChildRunners(role.Variables), playbook.Path)
			logger.Info("Finished Applying %s", host.ServerGroup)
		}
	}
}

func runTasks(role azad.Role, runners runners, rootPath string) error {
	for _, task := range role.Tasks {
		taskSchema, _ := communicator.GetTask(task.PluginName(), task.TaskName())
		for _, runner := range runners {
			runTaskParams := runTaskParams{
				task:       task,
				taskSchema: taskSchema,
				runner:     runner,
				rootPath:   rootPath,
				rolePath:   role.Path,
				user:       role.User,
			}
			runTask(runTaskParams)
		}
	}
	return nil
}
