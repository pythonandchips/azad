package runner

import (
	"github.com/pythonandchips/azad/conn"
)

// Config for azad
type Config struct {
	KeyFilePath string
	User        string
}

// SSHConfig extract config from main configuration
func (config Config) SSHConfig() conn.Config {
	return conn.Config{
		KeyFilePath: valueOrDefault(config.KeyFilePath, conn.DefaultSSHKeyPath()),
		User:        valueOrDefault(config.User, conn.DefaultUser()),
	}
}

func valueOrDefault(value, defaultValue string) string {
	if value == "" {
		return defaultValue
	}
	return value
}
