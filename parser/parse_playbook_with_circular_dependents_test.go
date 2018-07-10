package parser

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPlaybookWithCircularDependents(t *testing.T) {
	wd, _ := os.Getwd()
	filePath := filepath.Join(wd, "fixtures", "circular_dependents.az")
	env := map[string]string{
		"user":      "bruce_banner",
		"home_path": "/home/bruce_banner",
	}

	_, err := PlaybookFromFile(filePath, env)

	if err == nil {
		t.Fatalf("Expected but got none")
	}
	assert.Equal(t, err.Error(), "1 error occurred:\n\n* dependent Loop detected elasticsearch > java > elasticsearch")
}
