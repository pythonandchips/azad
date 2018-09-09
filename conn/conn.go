package conn

import (
	"bytes"
)

// Conn wraps connection, run and close of connection
type Conn interface {
	ConnectTo(string) error
	Run(Command) (Response, error)
	Close()
	Address() string
}

// NewConn create a now connection
var NewConn = newConn

func newConn(config Config) Conn {
	return &SSHConn{
		config: config,
	}
}

func newFakeConn(config Config) Conn {
	return &LoggerSSHConn{config: config}
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
