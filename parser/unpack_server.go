package parser

import (
	"fmt"

	"github.com/hashicorp/hcl2/hcl"
	"github.com/pythonandchips/azad/steps"
)

func unpackServer(block *hcl.Block) (steps.Step, error) {
	content, err := block.Body.Content(serverSchema())
	if err != nil {
		fmt.Println(err)
		return steps.ServerStep{}, err
	}
	return steps.ServerStep{
		Name:      block.Labels[0],
		Addresses: content.Attributes["addresses"],
	}, nil
}

func serverSchema() *hcl.BodySchema {
	return &hcl.BodySchema{
		Attributes: []hcl.AttributeSchema{
			{Name: "addresses", Required: true},
		},
	}
}
