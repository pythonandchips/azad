package conn

import (
	"bytes"
	"path/filepath"

	homedir "github.com/mitchellh/go-homedir"
)

type Conn interface {
	ConnectTo(string) error
	Run(Command) (Response, error)
	Close()
}

var NewConn = newConn

func newConn(config Config) Conn {
	return &SSHConn{
		config: config,
	}
}

func newFakeConn(config Config) Conn {
	return &LoggerSSHConn{}
}

// DefaultSSHKeyPath returns default path for ssh key
//
// $HOME/.ssh/id_rsa
func DefaultSSHKeyPath() string {
	home, _ := homedir.Dir()
	return filepath.Join(home, ".ssh", "id_rsa")
}

// Simulate switch connection to output to
// stdout instead of run on server
var Simulate = func() {
	NewConn = newFakeConn
}

type sshClient interface {
	NewSession() (sshSession, error)
	Close() error
}

type sshSession interface {
	setStdout(*bytes.Buffer)
	setStderr(*bytes.Buffer)
	Run(string) error
	Close() error
}
