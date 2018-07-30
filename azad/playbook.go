package azad

import (
	"fmt"

	"github.com/pythonandchips/azad/helpers/stringslice"
	"github.com/zclconf/go-cty/cty"
)

// Playbook describes servers and there roles
type Playbook struct {
	Inventories []Inventory
	Servers     []Server
	Hosts       []Host
	Variables   map[string]cty.Value
}

// LookupServer lookup server
func (playbook Playbook) LookupServer(hostName string) (Server, error) {
	for _, server := range playbook.Servers {
		if server.Group == hostName {
			return server, nil
		}
	}
	return Server{}, fmt.Errorf("server %s not found", hostName)
}

// AddAddressToServerByGroup adds an address to an existing server or adds
// a new server if one does not exist with the same group
// does not add the address if it already exists in the group
func (playbook *Playbook) AddAddressToServerByGroup(group, address string) {
	for i, server := range playbook.Servers {
		if server.Group == group {
			if stringslice.Exists(address, server.Addresses) {
				return
			}
			server.Addresses = append(server.Addresses, address)
			playbook.Servers[i] = server
			return
		}
	}
	server := Server{Group: group, Addresses: []string{address}}
	playbook.Servers = append(playbook.Servers, server)
}

// ListTasks list tasks
func (playbook Playbook) ListTasks() Tasks {
	tasks := Tasks{}
	for _, host := range playbook.Hosts {
		for _, role := range host.Roles {
			for _, task := range role.Tasks {
				tasks = append(tasks, task)
			}
		}
	}
	return tasks
}

// RequiredTasks list tasks required for playbook
func (playbook Playbook) RequiredTasks() Tasks {
	return playbook.ListTasks().UniqueTypes()
}
