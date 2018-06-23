package parser

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseRolesPerFolder(t *testing.T) {
	wd, _ := os.Getwd()
	filePath := filepath.Join(wd, "fixtures", "roles_per_folder", "basic.az")
	env := map[string]string{}

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
		if len(hosts) != 1 {
			t.Fatalf("Expected %d hosts but got %d", 1, len(hosts))
		}

		host := hosts[0]
		assert.Equal(t, host.ServerGroup, "tag_kibana_server")
		if len(host.Roles) != 1 {
			t.Fatalf("Expected host with %d roles but got %d", 1, len(host.Roles))
		}
		assert.Equal(t, host.Roles, []string{"ruby"})
	})

	t.Run("returns the roles configuration", func(t *testing.T) {
		roles := playbook.Roles
		if len(roles) != 3 {
			t.Fatalf("Expected %d roles but got %d", 3, len(roles))
		}
		assert.True(t, playbook.ContainsRole("ruby"), "Expected playbook to contain ruby role")
		assert.True(t, playbook.ContainsRole("security/firewall"), "expected playbook to contain security/firewall role")
		assert.True(t, playbook.ContainsRole("security/patches"), "expected playbook to contain security/patches role")

		t.Run("parses tasks for role", func(t *testing.T) {
			rubyRole, _ := playbook.FindRole("ruby")

			if len(rubyRole.Tasks) != 2 {
				t.Fatalf("Expected %d tasks but got %d", 2, len(rubyRole.Tasks))
			}
		})
	})
}
