package parser

import (
	"github.com/hashicorp/hcl2/hcl"
	"github.com/pythonandchips/azad/steps"
)

func unpackContext(block *hcl.Block) (steps.ContextContainer, error) {
	contextContainer := steps.ContextContainer{
		Name: block.Labels[0],
	}
	content, _ := block.Body.Content(contextSchema())
	attrs := content.Attributes
	for name, attr := range attrs {
		switch name {
		case "user":
			contextContainer.User = attr
		case "apply-to":
			contextContainer.ApplyTo = attr
		}
	}
	contextContainer.Steps, _ = unpackSteps(content.Blocks)
	return contextContainer, nil
}

func contextSchema() *hcl.BodySchema {
	return &hcl.BodySchema{
		Attributes: []hcl.AttributeSchema{
			{Name: "user"},
			{Name: "apply-to"},
		},
		Blocks: []hcl.BlockHeaderSchema{
			{Type: "variable", LabelNames: []string{"name"}},
			{Type: "input", LabelNames: []string{"type"}},
			{Type: "task", LabelNames: []string{"type", "name"}},
			{Type: "includes"},
		},
	}
}
