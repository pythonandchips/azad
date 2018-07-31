package parser

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/pythonandchips/azad/expect"
	"github.com/stretchr/testify/assert"
)

func TestParseRolesPerFolder(t *testing.T) {
	wd, _ := os.Getwd()
	filePath := filepath.Join(wd, "fixtures", "roles_per_folder")
	playbookPath := filepath.Join(filePath, "basic.az")
	env := map[string]string{}

	playbook, err := PlaybookFromFile(playbookPath, env)

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

		t.Run("host contains the roles for configuration", func(t *testing.T) {
			roles := host.Roles

			expect.EqualFatal(t, len(roles), 3)

			assert.Equal(t, roles[0].Name, "security/firewall")
			assert.Equal(t, roles[1].Name, "security/patches")
			assert.Equal(t, roles[2].Name, "ruby")

			t.Run("parse dependent roles", func(t *testing.T) {
				role := roles[0]
				rolePath := filepath.Join(filePath, "roles", "security", "firewall")
				assert.Equal(t, role.Path, rolePath)

				assert.Equal(t, len(role.Tasks), 1)
				assert.Equal(t, len(role.Variables), 1)
			})

			t.Run("parses tasks for role", func(t *testing.T) {
				rubyRole := roles[2]
				rolePath := filepath.Join(filePath, "roles", "ruby")
				assert.Equal(t, rubyRole.Path, rolePath)

				assert.Equal(t, len(rubyRole.Tasks), 3)
			})
		})
	})
}
