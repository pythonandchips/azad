package azad

import "fmt"

// Playbook describes servers and there roles
type Playbook struct {
	Servers []Server
	Hosts   []Host
	Roles   []Role
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
	tasks := Tasks{}
	for _, roleName := range roleNames {
		role, err := playbook.FindRole(roleName)
		if err != nil {
			return tasks, err
		}
		tasks = append(tasks, role.Tasks...)
	}
	return tasks, nil
}

// FindRole find role
func (playbook Playbook) FindRole(roleName string) (Role, error) {
	for _, role := range playbook.Roles {
		if role.Name == roleName {
			return role, nil
		}
	}
	return Role{}, fmt.Errorf("Role %s not found", roleName)
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
