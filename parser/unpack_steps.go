package parser

import (
	"github.com/hashicorp/hcl2/hcl"
	"github.com/pythonandchips/azad/steps"
)

func unpackSteps(blocks hcl.Blocks) (steps.StepList, error) {
	stepList := steps.StepList{}
	for _, block := range blocks {
		var step steps.Step
		switch block.Type {
		case "server":
			step, _ = unpackServer(block)
		case "inventory":
			step = steps.InventoryStep{
				Type: block.Labels[0],
				Body: block.Body,
			}
		case "variable":
			step, _ = unpackVariable(block)
		case "input":
			step = steps.InputStep{
				Type: block.Labels[0],
				Name: block.Labels[1],
				Body: block.Body,
			}
		case "task":
			step, _ = unpackTask(block)
		case "includes":
			step, _ = unpackIncludes(block)
		case "context":
			step, _ = unpackContext(block)
		case "role":
			// ignore role we will handle this seperatly and add to a role list
			continue
		default:
			continue
		}
		stepList = append(stepList, step)
	}
	return stepList, nil
}
