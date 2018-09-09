package parser

import (
	"github.com/hashicorp/hcl2/hcl"
	"github.com/pythonandchips/azad/steps"
)

func unpackVariable(block *hcl.Block) (steps.Step, error) {
	content, err := block.Body.Content(variableSchema())
	if err != nil {
		return steps.VariableStep{}, err
	}
	return steps.VariableStep{
		Name:  block.Labels[0],
		Type:  content.Attributes["type"],
		Value: content.Attributes["value"],
	}, nil
}

func variableSchema() *hcl.BodySchema {
	return &hcl.BodySchema{
		Attributes: []hcl.AttributeSchema{
			{Name: "type", Required: false},
			{Name: "value", Required: true},
		},
	}
}
