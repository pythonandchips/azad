package azad

import "github.com/pythonandchips/azad/conn"

// Config for azad
type Config struct {
	KeyFilePath string
}

// SSHConfig extract config from main configuration
func (config Config) SSHConfig() conn.Config {
	keyFilePath := config.KeyFilePath
	if keyFilePath == "" {
		keyFilePath = conn.DefaultSSHKeyPath()
	}
	return conn.Config{
		KeyFilePath: keyFilePath,
	}
}
