package azad

import (
	"fmt"
	"strings"

	"github.com/pythonandchips/azad/helpers/stringslice"
)

// Playbook describes servers and there roles
type Playbook struct {
	Inventories []Inventory
	Servers     []Server
	Hosts       []Host
	Roles       Roles
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
	for _, role := range playbook.Roles {
		for _, task := range role.Tasks {
			tasks = append(tasks, task)
		}
	}
	return tasks
}

// RequiredTasks list tasks required for playbook
func (playbook Playbook) RequiredTasks() Tasks {
	return playbook.ListTasks().UniqueTypes()
}

// TasksForRoles tasks for role
func (playbook Playbook) TasksForRoles(roleNames []string) (Tasks, error) {
	return playbook.tasksForRolesWithDependentCheck(roleNames, []string{})
}

func (playbook Playbook) tasksForRolesWithDependentCheck(roleNames []string, pastRuns []string) (Tasks, error) {
	tasks := Tasks{}
	for _, roleName := range roleNames {
		if stringslice.Exists(roleName, pastRuns) {
			pastRuns = append(pastRuns, roleName)
			path := strings.Join(pastRuns, " > ")
			return tasks, fmt.Errorf("dependent Loop detected %s", path)
		}
		pastRuns = append(pastRuns, roleName)
		role, err := playbook.FindRole(roleName)
		if err != nil {
			return tasks, err
		}
		dependentTasks, err := playbook.tasksForRolesWithDependentCheck(role.Dependents, pastRuns)
		if err != nil {
			return tasks, err
		}
		tasks = append(tasks, dependentTasks...)
		tasks = append(tasks, role.Tasks...)
		pastRuns = pastRuns[:len(pastRuns)-1]
	}
	return tasks, nil
}

// FindRole find role
func (playbook Playbook) FindRole(roleName string) (Role, error) {
	return playbook.Roles.FindRole(roleName)
}

// ContainsRole test if playbook contains a role by name
func (playbook Playbook) ContainsRole(roleName string) bool {
	for _, role := range playbook.Roles {
		if role.Name == roleName {
			return true
		}
	}
	return false
}
