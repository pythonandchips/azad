package azad

import (
	"testing"

	"github.com/pythonandchips/azad/conn"
	"github.com/stretchr/testify/assert"
)

func TestConfigSSHConfig(t *testing.T) {
	t.Run("returns the ssh config", func(t *testing.T) {
		config := Config{
			KeyFilePath: "keyFilePath",
		}
		connConfig := config.SSHConfig()

		assert.Equal(t, connConfig.KeyFilePath, config.KeyFilePath)
	})

	t.Run("default path if ssh key not set", func(t *testing.T) {
		config := Config{}
		connConfig := config.SSHConfig()

		assert.Equal(t, connConfig.KeyFilePath, conn.DefaultSSHKeyPath())
	})
}
