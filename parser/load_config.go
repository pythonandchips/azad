package parser

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/hashicorp/hcl2/gohcl"
	"github.com/hashicorp/hcl2/hcl"
	"github.com/hashicorp/hcl2/hcl/hclsyntax"
)

func loadConfigFile(filename string) (playbookDescription, error) {
	playbookDescription, err := parsePlaybookConfig(filename)
	if err != nil {
		return playbookDescription, err
	}
	roleDescriptions, err := parseRoles(filename)
	if err != nil {
		return playbookDescription, err
	}
	playbookDescription.Roles = append(playbookDescription.Roles, roleDescriptions...)
	return playbookDescription, err
}

func parsePlaybookConfig(filename string) (playbookDescription, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return playbookDescription{}, nil
	}
	playbook := playbookDescription{}
	file, parseErr := hclsyntax.ParseConfig(
		data,
		"config",
		hcl.Pos{Line: 1, Column: 1},
	)

	if parseErr.HasErrors() {
		return playbook, err
	}

	if parseErr = gohcl.DecodeBody(file.Body, nil, &playbook); parseErr.HasErrors() {
		return playbook, parseErr
	}

	return playbook, nil
}

func parseRoles(filename string) (roleDescriptions, error) {
	dirPath := filepath.Dir(filename)
	rolesFolderPath := filepath.Join(dirPath, "roles")
	roleDescriptions := roleDescriptions{}

	err := filepath.Walk(rolesFolderPath, func(path string, info os.FileInfo, err error) error {
		if path == rolesFolderPath {
			return nil
		}
		roleDescriptions, err = handleFilePath(roleDescriptions, rolesFolderPath, path, info)
		return err
	})

	return roleDescriptions, err
}

func handleFilePath(
	roleDescriptions roleDescriptions,
	rolesFolderPath, path string,
	info os.FileInfo,
) (roleDescriptions, error) {
	// Ignore roles folder, all roles must be contained in a folder
	if info.IsDir() {
		roleDescription, err := createRoleDescriptionFromFolder(rolesFolderPath, path)
		if err != nil {
			return roleDescriptions, nil
		}
		roleDescriptions = append(roleDescriptions, roleDescription)
		return roleDescriptions, nil
	}
	// only parse file with .az extension
	if filepath.Ext(path) != ".az" {
		return roleDescriptions, nil
	}
	err := createRoleFromConfig(path, rolesFolderPath, roleDescriptions)
	return roleDescriptions, err
}

func createRoleFromConfig(
	path, rolesFolderPath string,
	roleDescriptions roleDescriptions,
) error {
	dirPath := filepath.Dir(path)
	roleName, err := filepath.Rel(rolesFolderPath, dirPath)
	if err != nil {
		return err
	}
	role, err := roleDescriptions.RoleFor(roleName)
	if err != nil {
		return err
	}
	err = parseRoleConfig(path, &role)
	if err != nil {
		return fmt.Errorf("Failed to parse role %s: %s", roleName, err)
	}
	roleDescriptions.ReplaceWith(role)
	return nil
}

func parseRoleConfig(path string, role *roleDescription) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	file, parseErr := hclsyntax.ParseConfig(
		data,
		"config",
		hcl.Pos{Line: 1, Column: 1},
	)

	if parseErr.HasErrors() {
		return parseErr
	}
	if parseErr = gohcl.DecodeBody(file.Body, nil, role); parseErr.HasErrors() {
		return parseErr
	}
	return nil
}

func createRoleDescriptionFromFolder(rolesFolderPath, path string) (roleDescription, error) {
	globPath := filepath.Join(path, "*.az")
	azadFiles, _ := filepath.Glob(globPath)
	// Ignore folders that do not contain any azad file
	if len(azadFiles) == 0 {
		return roleDescription{}, fmt.Errorf("No azad files found")
	}
	roleName, err := filepath.Rel(rolesFolderPath, path)
	if err != nil {
		return roleDescription{}, fmt.Errorf("Unable to find relative path")
	}
	roleDescription := roleDescription{
		Name: roleName,
	}
	return roleDescription, nil
}
