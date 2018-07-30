package runner

import (
	"testing"

	"github.com/pythonandchips/azad/conn"
	"github.com/stretchr/testify/assert"
	"github.com/zclconf/go-cty/cty"
)

func RunnerEquals(t *testing.T, runner *runner, address string, variables map[string]cty.Value) {
	conn, _ := runner.Conn.(*fakeConn)
	assert.Equal(t, runner.Address, address)
	assert.Equal(t, runner.GlobalVariables, variables)
	assert.Equal(t, conn.ConnectedTo, address)
}

type fakeConn struct {
	ConnectedTo     string
	ConnectionError error
	Closed          bool
}

func (fakeConn *fakeConn) ConnectTo(ip string) error {
	fakeConn.ConnectedTo = ip
	return fakeConn.ConnectionError
}

func (fakeConn *fakeConn) Run(conn.Command) (conn.Response, error) {
	return conn.CommandResponse{}, nil
}

func (fakeConn *fakeConn) Close() {
	fakeConn.Closed = true
}
