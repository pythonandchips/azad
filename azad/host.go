package azad

import "github.com/zclconf/go-cty/cty"

// Host name and roles of server
type Host struct {
	ServerGroup string
	Roles       []Role
	Variables   map[string]cty.Value
}
