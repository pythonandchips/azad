package conn

import (
	"os/user"
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
	currentUser = func() (*user.User, error) {
		return &user.User{
			Username: "bruce_springsteen",
		}, nil
	}

	assert.Equal(t, DefaultUser(), "bruce_springsteen")
}
