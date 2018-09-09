package plugin

import "github.com/pythonandchips/azad/conn"

type FakeConn struct {
	response conn.Response
	command  conn.Command
}

func (fakeConn *FakeConn) ConnectTo(string) error {
	return nil
}

func (fakeConn *FakeConn) Run(command conn.Command) (conn.Response, error) {
	fakeConn.command = command
	return fakeConn.response, nil
}

func (fakeConn *FakeConn) Close() {}

func (fakeConn *FakeConn) Address() string {
	return ""
}

type FakeResponse struct {
	stdout string
	stderr string
}

func (fakeResponse FakeResponse) Stdout() string {
	return fakeResponse.stdout
}

func (fakeResponse FakeResponse) Stderr() string {
	return fakeResponse.stderr
}
