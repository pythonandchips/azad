package parser

import (
	"fmt"
	"io/ioutil"

	"github.com/hashicorp/hcl2/hcl"
	"github.com/hashicorp/hcl2/hcl/hclsyntax"
	"github.com/pythonandchips/azad/steps"
)

// PlaybookSteps parse out the steps for a playbook
func PlaybookSteps(path string) (steps.PlaybookSteps, error) {
	playbookSteps := steps.PlaybookSteps{}
	body, err := parseFile(path)
	if err != nil {
		return playbookSteps, err
	}
	content, diag := body.Content(playbookSchema())
	playbookSteps.StepList, _ = unpackSteps(content.Blocks)
	playbookSteps.RoleList, _ = unpackRoles(path, content.Blocks)
	if diag.HasErrors() {
		return playbookSteps, fmt.Errorf(diag.Error())
	}
	return playbookSteps, nil
}

func parseFile(path string) (hcl.Body, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return hcl.EmptyBody(), err
	}
	file, parseErr := hclsyntax.ParseConfig(
		data,
		"config",
		hcl.Pos{Line: 1, Column: 1},
	)
	if parseErr.HasErrors() {
		return hcl.EmptyBody(), fmt.Errorf(parseErr.Error())
	}
	return file.Body, nil
}

func playbookSchema() *hcl.BodySchema {
	return &hcl.BodySchema{
		Blocks: []hcl.BlockHeaderSchema{
			{Type: "server", LabelNames: []string{"name"}},
			{Type: "inventory", LabelNames: []string{"type"}},
			{Type: "variable", LabelNames: []string{"name"}},
			{Type: "input", LabelNames: []string{"type", "name"}},
			{Type: "context", LabelNames: []string{"name"}},
			{Type: "role", LabelNames: []string{"name"}},
		},
	}
}
