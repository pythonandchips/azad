package main

import (
	"testing"

	"github.com/pythonandchips/azad/azad"
	"github.com/pythonandchips/azad/conn"
	"github.com/pythonandchips/azad/logger"
	"github.com/pythonandchips/azad/runner"
	"github.com/stretchr/testify/assert"
)

func TestRunPlaybook(t *testing.T) {
	logger.StubLogger()
	var specifiedPlaybookFilePath string
	var specifiedConfig azad.Config
	simulation := false
	runner.RunPlaybook = func(playbookFilePath string, config azad.Config) {
		specifiedPlaybookFilePath = playbookFilePath
		specifiedConfig = config
	}
	conn.Simulate = func() {
		simulation = true
	}
	t.Run("passes the correct playbook path to the runner", func(t *testing.T) {
		fakeContext := FakeContext{
			args: []string{"test_playbook.az"},
			strings: map[string]string{
				"key":  "path_to_key",
				"user": "admin",
			},
		}
		runPlaybook(fakeContext)
		assert.Equal(t, specifiedPlaybookFilePath, "test_playbook.az")
		assert.Equal(t, specifiedConfig.KeyFilePath, "path_to_key")
		assert.Equal(t, specifiedConfig.User, "admin")
		assert.False(t, simulation)
	})
	t.Run("uses the default playbook if one is not specified", func(t *testing.T) {
		fakeContext := FakeContext{
			args: []string{},
			strings: map[string]string{
				"key": "path_to_key",
			},
		}
		runPlaybook(fakeContext)
		assert.Equal(t, specifiedPlaybookFilePath, "./playbook.az")
	})
	t.Run("set conn to simulate when simulate switch set", func(t *testing.T) {
		fakeContext := FakeContext{
			args: []string{},
			bools: map[string]bool{
				"simulate": true,
			},
		}
		runPlaybook(fakeContext)
		assert.True(t, simulation)
	})
}
