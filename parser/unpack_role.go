package parser

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/hcl2/hcl"
	"github.com/pythonandchips/azad/steps"
)

func unpackRoles(path string, blocks hcl.Blocks) ([]steps.RoleContainer, error) {
	roleContainers := []steps.RoleContainer{}
	absPath, _ := filepath.Abs(filepath.Dir(path))
	for _, block := range blocks {
		if block.Type != "role" {
			continue
		}
		name := block.Labels[0]
		content, _ := block.Body.Content(roleSchema())
		role, _ := unpackRole(name, content, absPath)
		roleContainers = append(roleContainers, role)
	}
	roleContainersFromfile, _ := loadRolesFromFolders(path)
	roleContainers = append(roleContainers, roleContainersFromfile...)

	return roleContainers, nil
}

func loadRolesFromFolders(path string) ([]steps.RoleContainer, error) {
	dirPath := filepath.Dir(path)
	rolesFolderPath := filepath.Join(dirPath, "roles")
	roleContainers := []steps.RoleContainer{}
	err := filepath.Walk(rolesFolderPath, func(path string, info os.FileInfo, err error) error {
		if path == rolesFolderPath {
			return nil
		}
		roleContainers, err = handleFilePath(roleContainers, rolesFolderPath, path, info)
		return err
	})
	return roleContainers, err
}

func handleFilePath(
	roleContainers []steps.RoleContainer,
	rolesFolderPath, path string,
	info os.FileInfo,
) ([]steps.RoleContainer, error) {
	if info.IsDir() {
		return roleContainers, nil
	}
	// only parse file with .az extension
	if filepath.Ext(path) != ".az" {
		return roleContainers, nil
	}
	body, _ := parseFile(path)
	content, _ := body.Content(roleSchema())
	rolePath, _ := filepath.Rel(rolesFolderPath, path)
	roleName := strings.TrimSuffix(rolePath, ".az")
	absPath, _ := filepath.Abs(filepath.Dir(path))
	if strings.HasSuffix(roleName, "/main") {
		roleName = strings.Replace(roleName, "/main", "", -1)
	}
	roleContainer, _ := unpackRole(roleName, content, absPath)
	roleContainers = append(roleContainers, roleContainer)

	return roleContainers, nil
}

func roleSchema() *hcl.BodySchema {
	return &hcl.BodySchema{
		Attributes: []hcl.AttributeSchema{
			{Name: "user"},
		},
		Blocks: []hcl.BlockHeaderSchema{
			{Type: "variable", LabelNames: []string{"name"}},
			{Type: "input", LabelNames: []string{"type", "name"}},
			{Type: "task", LabelNames: []string{"type", "name"}},
			{Type: "includes"},
		},
	}
}

func unpackRole(name string, content *hcl.BodyContent, roleFile string) (steps.RoleContainer, error) {
	roleContainer := steps.RoleContainer{
		Name: name,
		File: roleFile,
	}
	attrs := content.Attributes
	for name, attr := range attrs {
		switch name {
		case "user":
			roleContainer.User = attr
		}
	}
	roleContainer.Steps, _ = unpackSteps(content.Blocks)
	return roleContainer, nil
}
