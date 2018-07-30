package runner

import (
	multierror "github.com/hashicorp/go-multierror"
	"github.com/pythonandchips/azad/azad"
	"github.com/pythonandchips/azad/communicator"
)

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
