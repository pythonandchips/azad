package parser

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/pythonandchips/azad/expect"
	"github.com/stretchr/testify/assert"
)

func TestPlaybookFromFileBasic(t *testing.T) {
	wd, _ := os.Getwd()
	filePath := filepath.Join(wd, "fixtures", "basic.az")
	env := map[string]string{
		"user":      "bruce_banner",
		"home_path": "/home/bruce_banner",
	}

	playbook, err := PlaybookFromFile(filePath, env)

	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	t.Run("returns the server with addresses", func(t *testing.T) {
		servers := playbook.Servers
		if len(servers) != 1 {
			t.Fatalf("Expected %d servers but got %d", 1, len(servers))
		}

		server := servers[0]
		assert.Equal(t, server.Group, "tag_kibana_server")
		if len(server.Addresses) != 1 {
			t.Fatalf("Expected server with %d addresses but had %d", 1, len(server.Addresses))
		}

		address := server.Addresses[0]
		assert.Equal(t, address, "10.0.0.1")
	})

	t.Run("returns the host configuration", func(t *testing.T) {
		hosts := playbook.Hosts
		expect.EqualFatal(t, len(hosts), 1)

		host := hosts[0]
		assert.Equal(t, host.ServerGroup, "tag_kibana_server")

		expect.EqualFatal(t, len(host.Roles), 1)

		expect.EqualFatal(t, len(host.Variables), 2)
		assert.Equal(t, host.Variables["install_path"].AsString(), "/opt/installer")

		t.Run("returns the roles configuration", func(t *testing.T) {
			roles := host.Roles

			expect.EqualFatal(t, len(roles), 1)

			role := roles[0]
			assert.Equal(t, role.Name, "elasticsearch")
			assert.Equal(t, role.User, "root")
			expect.EqualFatal(t, len(role.Variables), 1)
			assert.Equal(t, role.Variables["ruby_install_path"].AsString(), "/opt/installer/ruby")

			t.Run("parses the tasks for the role", func(t *testing.T) {
				tasks := role.Tasks
				if len(tasks) != 2 {
					t.Fatalf("Expected %d tasks for %s role but got %d", 1, role.Name, len(tasks))
				}

				task := role.Tasks[0]
				assert.Equal(t, task.Type, "stat")
				assert.Equal(t, task.Name, "ruby-exists")

				assert.Equal(t, task.Attributes["path"].Name, "path")
				t.Run("parse task with a conditional", func(t *testing.T) {
					conditionalTask := role.Tasks[1]
					assert.NotNil(t, conditionalTask.Condition)
				})
			})
		})
	})
}
