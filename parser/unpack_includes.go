package parser

import (
	"github.com/hashicorp/hcl2/hcl"
	"github.com/pythonandchips/azad/steps"
)

func unpackIncludes(block *hcl.Block) (steps.Step, error) {
	content, err := block.Body.Content(includeSchema())
	if err != nil {
		return steps.IncludesStep{}, err
	}
	return steps.IncludesStep{
		Roles: content.Attributes["roles"],
	}, nil
}

func includeSchema() *hcl.BodySchema {
	return &hcl.BodySchema{
		Attributes: []hcl.AttributeSchema{
			{Name: "roles", Required: true},
		},
	}
}
