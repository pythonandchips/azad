package parser

import (
	"io/ioutil"

	"github.com/hashicorp/hcl2/gohcl"
	"github.com/hashicorp/hcl2/hcl"
	"github.com/hashicorp/hcl2/hcl/hclsyntax"
)

func loadConfigFile(fileName string) (playbookDescription, error) {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return playbookDescription{}, err
	}
	return parseConfig(data)
}

func parseConfig(playbookConfig []byte) (playbookDescription, error) {
	playbook := playbookDescription{}
	file, err := hclsyntax.ParseConfig(
		playbookConfig,
		"config",
		hcl.Pos{Line: 1, Column: 1},
	)

	if err.HasErrors() {
		return playbook, err
	}

	if err = gohcl.DecodeBody(file.Body, nil, &playbook); err.HasErrors() {
		return playbook, err
	}

	return playbook, nil
}
