package parser

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseRoleConfig(t *testing.T) {
	wd, _ := os.Getwd()
	filePath := filepath.Join(
		wd, "fixtures", "roles_per_folder",
		"roles", "ruby", "main.az",
	)
	roleFile, err := parseRoleConfig(filePath)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	assert.Equal(t, roleFile.name, "main")
	roleDescription := roleFile.roleDescription
	assert.Equal(t, roleDescription.Name, "main")
	roleSchema := roleSchema()
	content, err := roleDescription.Config.Content(&roleSchema)
	if len(content.Blocks) != 3 {
		t.Fatalf("expected 3 content blocks but got %d", len(content.Blocks))
	}
	assert.Equal(t, content.Blocks[0].Type, "task")
	assert.Equal(t, content.Blocks[1].Type, "include")
	assert.Equal(t, content.Blocks[2].Type, "task")
}
