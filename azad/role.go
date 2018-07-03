package azad

import "fmt"

// Roles collection of role's
type Roles []Role

// FindRole find role in collection
func (roles Roles) FindRole(roleName string) (Role, error) {
	for _, role := range roles {
		if role.Name == roleName {
			return role, nil
		}
	}
	return Role{}, fmt.Errorf("Role %s not found", roleName)
}

// Role list of task to be applied to host
type Role struct {
	Name       string
	Dependents []string
	Tasks      Tasks
}
