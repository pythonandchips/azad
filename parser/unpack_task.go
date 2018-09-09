package parser

import (
	"github.com/hashicorp/hcl2/hcl"
	"github.com/pythonandchips/azad/steps"
)

func unpackTask(block *hcl.Block) (steps.TaskStep, error) {
	content, body, err := block.Body.PartialContent(taskSchema())
	return steps.TaskStep{
		Type:      block.Labels[0],
		Label:     block.Labels[1],
		User:      content.Attributes["user"],
		Condition: content.Attributes["condition"],
		Debug:     content.Attributes["debug"],
		Body:      body,
	}, err
}

func taskSchema() *hcl.BodySchema {
	return &hcl.BodySchema{
		Attributes: []hcl.AttributeSchema{
			{Name: "user", Required: false},
			{Name: "condition", Required: false},
			{Name: "debug", Required: false},
		},
	}
}
