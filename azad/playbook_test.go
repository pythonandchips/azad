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

func TestPlaybookAddAddressToServerGroup(t *testing.T) {
	t.Run("when group does not exist", func(t *testing.T) {
		playbook := Playbook{}
		playbook.AddAddressToServerByGroup("new_group", "10.0.0.1")

		if len(playbook.Servers) != 1 {
			t.Fatalf("expected %d servers but got %d", 1, len(playbook.Servers))
		}
		server := playbook.Servers[0]
		assert.Equal(t, server.Group, "new_group")

		if len(server.Addresses) != 1 {
			t.Fatalf("expected server to have %d addresses but got %d", 1, len(server.Addresses))
		}
		assert.Equal(t, server.Addresses[0], "10.0.0.1")
	})
	t.Run("when group already exists", func(t *testing.T) {
		playbook := Playbook{
			Servers: []Server{
				{Group: "existing_group", Addresses: []string{"10.0.0.1"}},
			},
		}
		playbook.AddAddressToServerByGroup("existing_group", "10.0.0.2")
		if len(playbook.Servers) != 1 {
			t.Fatalf("expected %d servers but got %d", 1, len(playbook.Servers))
		}

		server := playbook.Servers[0]
		if len(server.Addresses) != 2 {
			t.Fatalf("expected server to have %d addresses but got %d", 2, len(server.Addresses))
		}
		assert.Equal(t, server.Addresses[1], "10.0.0.2")
	})
	t.Run("when group and address already exist", func(t *testing.T) {
		playbook := Playbook{
			Servers: []Server{
				{Group: "existing_group", Addresses: []string{"10.0.0.1"}},
			},
		}
		playbook.AddAddressToServerByGroup("existing_group", "10.0.0.1")
		if len(playbook.Servers) != 1 {
			t.Fatalf("expected %d servers but got %d", 1, len(playbook.Servers))
		}

		server := playbook.Servers[0]
		if len(server.Addresses) != 1 {
			t.Fatalf("expected server to have %d addresses but got %d", 1, len(server.Addresses))
		}
	})
}
