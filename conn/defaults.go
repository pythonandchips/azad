package conn

import (
	"os/user"
	"path/filepath"

	homedir "github.com/mitchellh/go-homedir"
)

// DefaultSSHKeyPath returns default path for ssh key
//
// $HOME/.ssh/id_rsa
func DefaultSSHKeyPath() string {
	home, _ := homedir.Dir()
	return filepath.Join(home, ".ssh", "id_rsa")
}

// DefaultUser returns default user
//
// root
func DefaultUser() string {
	user, _ := currentUser()
	return user.Username
}

var currentUser = func() (*user.User, error) {
	return user.Current()
}
