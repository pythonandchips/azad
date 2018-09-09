package parser

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/pythonandchips/azad/expect"
	"github.com/pythonandchips/azad/steps"
)

func TestMultipleRoles(t *testing.T) {
	wd, _ := os.Getwd()
	filePath := filepath.Join(wd, "fixtures", "roles_per_folder")
	playbookPath := filepath.Join(filePath, "basic.az")

	playbookSteps, err := PlaybookSteps(playbookPath)

	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	expect.EqualFatal(t, len(playbookSteps.RoleList), 6)

	roleListContains(t, playbookSteps.RoleList, "ruby")
	roleListContains(t, playbookSteps.RoleList, "ruby/install")
	roleListContains(t, playbookSteps.RoleList, "ruby/update")
	roleListContains(t, playbookSteps.RoleList, "ruby/dependencies")
	roleListContains(t, playbookSteps.RoleList, "security/firewall")
	roleListContains(t, playbookSteps.RoleList, "security/patches")
}

func roleListContains(t *testing.T, roleList []steps.RoleContainer, name string) {
	for _, roleContainer := range roleList {
		if roleContainer.Name == name {
			return
		}
	}
	t.Errorf("role %s not found in list", name)
}
