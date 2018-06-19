package parser

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPlaybookFromFileInventory(t *testing.T) {
	wd, _ := os.Getwd()
	filePath := filepath.Join(wd, "fixtures", "inventory.az")
	env := map[string]string{
		"user":      "bruce_banner",
		"home_path": "/home/bruce_banner",
	}

	playbook, err := PlaybookFromFile(filePath, env)

	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	t.Run("returns the inventory", func(t *testing.T) {
		inventories := playbook.Inventories
		if len(inventories) != 1 {
			t.Fatalf("Expected %d inventory but got %d", 1, len(inventories))
		}

		inventory := inventories[0]
		assert.Equal(t, inventory.Name, "aws.ec2")

		attributes := inventory.Attributes
		assert.Equal(t, attributes["access_key_id"], "access_key_id")
		assert.Equal(t, attributes["secret_key"], "secret_key")
	})
}
