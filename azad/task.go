package azad

import (
	"strings"

	"github.com/hashicorp/hcl2/hcl"
)

// Tasks list of tasks
type Tasks []Task

// ContainsType is task in list
func (tasks Tasks) ContainsType(task Task) bool {
	for _, taskInArray := range tasks {
		if taskInArray.Type == task.Type {
			return true
		}
	}
	return false
}

// UniqueTypes Tasks
func (tasks Tasks) UniqueTypes() Tasks {
	uniqueTasks := Tasks{}
	for _, task := range tasks {
		if uniqueTasks.ContainsType(task) {
			continue
		}
		uniqueTasks = append(uniqueTasks, task)
	}
	return uniqueTasks
}

// Task run command via ssh
type Task struct {
	Type       string
	Name       string
	Attributes map[string]*hcl.Attribute
}

// PluginName name of plugin
func (task Task) PluginName() string {
	parts := strings.Split(task.Type, ".")
	if len(parts) == 2 {
		return parts[0]
	}
	return "core"
}

// TaskName name of task
func (task Task) TaskName() string {
	parts := strings.Split(task.Type, ".")
	if len(parts) == 2 {
		return parts[1]
	}
	return parts[0]
}
