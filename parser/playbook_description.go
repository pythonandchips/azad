package parser

import (
	"github.com/hashicorp/hcl2/hcl"
)

type playbookDescription struct {
	Inventories           []inventoryDescription `hcl:"inventory,block"`
	Servers               []serverDescription    `hcl:"server,block"`
	Hosts                 []hostDescription      `hcl:"host,block"`
	Variables             []variableDescription  `hcl:"variable,block"`
	Roles                 []roleDescription      `hcl:"role,block"`
	roleDescriptionGroups []roleDescriptionGroup
}

type inventoryDescription struct {
	Name   string   `hcl:",label"`
	Config hcl.Body `hcl:",remain"`
}

type serverDescription struct {
	Group     string   `hcl:",label"`
	Addresses []string `hcl:"addresses"`
	Config    hcl.Body `hcl:",remain"`
}

// Variable value
type variableDescription struct {
	Name   string   `hcl:",label"`
	Config hcl.Body `hcl:",remain"`
}
