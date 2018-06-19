package conn

import (
	"testing"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/stretchr/testify/assert"
)

func TestDefaultSSHKeyPath(t *testing.T) {
	defaultSSHKeyPath := DefaultSSHKeyPath()
	home, _ := homedir.Dir()
	assert.Equal(t, defaultSSHKeyPath, home+"/.ssh/id_rsa")
}

func TestDefaultUser(t *testing.T) {
	assert.Equal(t, DefaultUser(), "root")
}
