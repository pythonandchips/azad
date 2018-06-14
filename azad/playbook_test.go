package azad

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPlaybookListRequirements(t *testing.T) {
	playbook := Playbook{
		Roles: []Role{
			{
				Name: "elastic_search",
				Tasks: []Task{
					{Type: "bash", Name: "list files"},
					{Type: "debian.apt-get", Name: "list files"},
					{Type: "credstash.lookup", Name: "list files"},
					{Type: "debian.apt-get", Name: "list files"},
					{Type: "debian.chmod", Name: "list files"},
				},
			},
		},
	}

	t.Run("list required plugins plugins", func(t *testing.T) {
		tasks := playbook.RequiredTasks()

		if len(tasks) != 4 {
			t.Errorf("Expected 3 tasks but got %d", len(tasks))
			t.FailNow()
		}

		assert.Equal(t, tasks[0].Type, "bash")
		assert.Equal(t, tasks[1].Type, "debian.apt-get")
		assert.Equal(t, tasks[2].Type, "credstash.lookup")
		assert.Equal(t, tasks[3].Type, "debian.chmod")
	})
}

func TestPlaybookTasksForRoles(t *testing.T) {
	playbook := Playbook{
		Roles: []Role{
			{
				Name: "install_ruby",
				Tasks: []Task{
					{Type: "debian.apt-get", Name: "Install ruby"},
				},
			},
			{
				Name: "elastic_search",
				Tasks: []Task{
					{Type: "bash", Name: "list files"},
					{Type: "credstash.lookup", Name: "list files"},
				},
			},
		},
	}
	t.Run("list tasks for given roles", func(t *testing.T) {
		roleNames := []string{"install_ruby", "elastic_search"}
		tasks, _ := playbook.TasksForRoles(roleNames)

		if len(tasks) != 3 {
			t.Errorf("Expected 3 tasks but got %d", len(tasks))
			t.FailNow()
		}

		assert.Equal(t, tasks[0].Type, "debian.apt-get")
		assert.Equal(t, tasks[1].Type, "bash")
		assert.Equal(t, tasks[2].Type, "credstash.lookup")
	})
}
