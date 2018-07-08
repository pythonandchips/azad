package parser

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/hcl2/gohcl"
	"github.com/hashicorp/hcl2/hcl"
	"github.com/hashicorp/hcl2/hcl/hclsyntax"
)

func loadConfigFile(filename string) (playbookDescription, error) {
	playbookDescription, err := parsePlaybookConfig(filename)
	if err != nil {
		return playbookDescription, err
	}
	roleDescriptionGroups, err := parseRoles(filename)
	if err != nil {
		return playbookDescription, err
	}
	playbookRoleDescriptions := createRoleDescriptionGroupRoleDescriptions(playbookDescription.Roles)
	roleDescriptionGroups = append(roleDescriptionGroups, playbookRoleDescriptions...)
	playbookDescription.roleDescriptionGroups = roleDescriptionGroups

	return playbookDescription, err
}

func createRoleDescriptionGroupRoleDescriptions(roleDescriptions roleDescriptions) roleDescriptionGroups {
	roleDescriptionGroups := roleDescriptionGroups{}
	for _, roleDescription := range roleDescriptions {
		roleDescriptionGroup := roleDescriptionGroup{
			name: roleDescription.Name,
			files: []roleFile{
				{name: "main", roleDescription: roleDescription},
			},
		}
		roleDescriptionGroups = append(roleDescriptionGroups, roleDescriptionGroup)
	}
	return roleDescriptionGroups
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

func parseRoles(filename string) (roleDescriptionGroups, error) {
	dirPath := filepath.Dir(filename)
	rolesFolderPath := filepath.Join(dirPath, "roles")
	roleDescriptionGroups := roleDescriptionGroups{}

	err := filepath.Walk(rolesFolderPath, func(path string, info os.FileInfo, err error) error {
		if path == rolesFolderPath {
			return nil
		}
		roleDescriptionGroups, err = handleFilePath(roleDescriptionGroups, rolesFolderPath, path, info)
		return err
	})

	return roleDescriptionGroups, err
}

func handleFilePath(
	roleDescriptionGroups roleDescriptionGroups,
	rolesFolderPath, path string,
	info os.FileInfo,
) (roleDescriptionGroups, error) {
	// Ignore roles folder, all roles must be contained in a folder
	if info.IsDir() {
		roleDescriptionGroup, err := createRoleDescriptionGroupFromFolder(rolesFolderPath, path)
		if err != nil {
			return roleDescriptionGroups, nil
		}
		roleDescriptionGroups = append(roleDescriptionGroups, roleDescriptionGroup)
		return roleDescriptionGroups, nil
	}
	// only parse file with .az extension
	if filepath.Ext(path) != ".az" {
		return roleDescriptionGroups, nil
	}
	return createRoleFromConfig(path, rolesFolderPath, roleDescriptionGroups)
}

func createRoleFromConfig(
	path, rolesFolderPath string,
	roleDescriptionGroups roleDescriptionGroups,
) (roleDescriptionGroups, error) {
	dirPath := filepath.Dir(path)
	roleName, err := filepath.Rel(rolesFolderPath, dirPath)
	if err != nil {
		return roleDescriptionGroups, err
	}
	roleDescriptionGroup, err := roleDescriptionGroups.RoleFor(roleName)
	if err != nil {
		return roleDescriptionGroups, err
	}
	roleFile, err := parseRoleConfig(path)
	if err != nil {
		return roleDescriptionGroups, fmt.Errorf("Failed to parse role %s: %s", roleName, err)
	}
	roleDescriptionGroup.files = append(roleDescriptionGroup.files, roleFile)
	roleDescriptionGroups.ReplaceWith(roleDescriptionGroup)
	return roleDescriptionGroups, nil
}

func parseRoleConfig(path string) (roleFile, error) {
	name := filepath.Base(path)
	name = strings.Replace(name, ".az", "", -1)
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return roleFile{}, err
	}
	file, parseErr := hclsyntax.ParseConfig(
		data,
		"config",
		hcl.Pos{Line: 1, Column: 1},
	)
	if parseErr.HasErrors() {
		return roleFile{}, parseErr
	}
	roleDescription := roleDescription{Name: name}
	parseErr = gohcl.DecodeBody(file.Body, nil, &roleDescription)
	if parseErr.HasErrors() {
		return roleFile{}, parseErr
	}
	roleFile := roleFile{name: name, roleDescription: roleDescription}
	return roleFile, nil
}

func createRoleDescriptionGroupFromFolder(rolesFolderPath, path string) (roleDescriptionGroup, error) {
	globPath := filepath.Join(path, "*.az")
	azadFiles, _ := filepath.Glob(globPath)
	// Ignore folders that do not contain any azad file
	if len(azadFiles) == 0 {
		return roleDescriptionGroup{}, fmt.Errorf("No azad files found")
	}
	roleName, err := filepath.Rel(rolesFolderPath, path)
	if err != nil {
		return roleDescriptionGroup{}, fmt.Errorf("Unable to find relative path")
	}
	roleDescriptionGroup := roleDescriptionGroup{
		name: roleName,
	}
	return roleDescriptionGroup, nil
}
