package parser

import "github.com/hashicorp/hcl2/hcl"

type playbookDescription struct {
	Servers   []serverDescription   `hcl:"server,block"`
	Hosts     []hostDescription     `hcl:"host,block"`
	Variables []variableDescription `hcl:"variable,block"`
	Roles     []roleDescription     `hcl:"role,block"`
}

type serverDescription struct {
	Group     string   `hcl:",label"`
	Addresses []string `hcl:"addresses"`
	Config    hcl.Body `hcl:",remain"`
}

// Variable value
type variableDescription struct {
	Name    string `hcl:",label"`
	Default string `hcl:"default"`
}

// Role list of task to be applied to host
type roleDescription struct {
	Name   string            `hcl:",label"`
	Tasks  []taskDescription `hcl:"task,block"`
	Config hcl.Body          `hcl:",remain"`
}

// Task run command via ssh
type taskDescription struct {
	Type   string   `hcl:"type,label"`
	Name   string   `hcl:"name,label"`
	Config hcl.Body `hcl:",remain"`
}
