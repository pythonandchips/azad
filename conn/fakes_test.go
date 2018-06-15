package conn

import "bytes"

type FakeSSHClient struct {
	session *FakeSSHSession
}

func (fakeSSHClient FakeSSHClient) NewSession() (sshSession, error) {
	return fakeSSHClient.session, nil
}

func (fakeSSHClient FakeSSHClient) Close() error {
	return nil
}

type FakeSSHSession struct {
	commands []string
	stdout   string
	stderr   string
}

func (fakeSSHSession *FakeSSHSession) setStdout(stdout *bytes.Buffer) {
	stdout.Write([]byte(fakeSSHSession.stdout))
}
func (fakeSSHSession *FakeSSHSession) setStderr(stderr *bytes.Buffer) {
	stderr.Write([]byte(fakeSSHSession.stderr))
}
func (fakeSSHSession *FakeSSHSession) Run(command string) error {
	fakeSSHSession.commands = append(fakeSSHSession.commands, command)
	return nil
}
func (fakeSSHSession *FakeSSHSession) Close() error {
	return nil
}
func (fakeSSHSession *FakeSSHSession) Clear() {
	fakeSSHSession.commands = []string{}
}
