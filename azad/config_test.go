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
			User:        "ssh-user",
		}
		connConfig := config.SSHConfig()

		assert.Equal(t, connConfig.KeyFilePath, config.KeyFilePath)
		assert.Equal(t, connConfig.User, config.User)
	})

	t.Run("returns defaults for config", func(t *testing.T) {
		config := Config{}
		connConfig := config.SSHConfig()

		assert.Equal(t, connConfig.KeyFilePath, conn.DefaultSSHKeyPath())
		assert.Equal(t, connConfig.User, conn.DefaultUser())
	})
}
