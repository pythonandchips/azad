package conn

import "github.com/pythonandchips/azad/logger"

type FakeSSHConn struct {
	ConnectedTo string
	Commands    []Command
	closed      bool
}

func (fakeSSHConn *FakeSSHConn) ConnectTo(hostName string) error {
	fakeSSHConn.ConnectedTo = hostName
	return nil
}

func (fakeSSHConn *FakeSSHConn) Run(command Command) error {
	fakeSSHConn.Commands = append(fakeSSHConn.Commands, command)
	logger.Debug("Running on %s", fakeSSHConn.ConnectedTo)
	logger.Debug(command.generateFile())
	return nil
}

func (fakeSSHConn *FakeSSHConn) Close() {
	logger.Debug("Closing connection to %s", fakeSSHConn.ConnectedTo)
	fakeSSHConn.closed = true
}
