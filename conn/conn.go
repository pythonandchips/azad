package conn

import "bytes"

type Conn interface {
	ConnectTo(string) error
	Run(Command) (Response, error)
	Close()
}

var NewConn = newConn

func newConn() Conn {
	return &SSHConn{}
}

func newFakeConn() Conn {
	return &LoggerSSHConn{}
}

func Simulate() {
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
