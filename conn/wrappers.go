package conn

import (
	"bytes"

	"golang.org/x/crypto/ssh"
)

type sshClientWrapper struct {
	*ssh.Client
}

func (sshClientWrapper sshClientWrapper) NewSession() (sshSession, error) {
	session, err := sshClientWrapper.Client.NewSession()
	return sshSessionWrapper{session}, err
}

type sshSessionWrapper struct {
	*ssh.Session
}

func (sshSessionWrapper sshSessionWrapper) setStdout(stdout *bytes.Buffer) {
	sshSessionWrapper.Stdout = stdout
}

func (sshSessionWrapper sshSessionWrapper) setStderr(stderr *bytes.Buffer) {
	sshSessionWrapper.Stderr = stderr
}

func (sshSessionWrapper sshSessionWrapper) Run(command string) error {
	return sshSessionWrapper.Run(command)
}
