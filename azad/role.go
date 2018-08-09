package azad

import (
	"fmt"

	"github.com/zclconf/go-cty/cty"
)

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
	Variables  map[string]cty.Value
	Path       string
	User       string
}
