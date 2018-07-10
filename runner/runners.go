package runner

import (
	multierror "github.com/hashicorp/go-multierror"
	"github.com/pythonandchips/azad/conn"
	"github.com/zclconf/go-cty/cty"
)

type runners []*runner

var newConn = func(config conn.Config) conn.Conn {
	return conn.NewConn(config)
}

func createRunners(addresses []string, variables map[string]cty.Value) (runners, error) {
	runners := runners{}
	errors := &multierror.Error{}
	for _, address := range addresses {
		connection := newConn(config.SSHConfig())
		err := connection.ConnectTo(address)
		if err != nil {
			errors = multierror.Append(errors, err)
			break
		}
		runner := newRunner(connection, address, variables)
		runners = append(runners, runner)
	}
	return runners, errors.ErrorOrNil()
}

func (runners runners) Close() {
	for _, runner := range runners {
		runner.Conn.Close()
	}
}

func (runners runners) ChildRunners(variables map[string]cty.Value) runners {
	childRunners := []*runner{}
	for _, runner := range runners {
		childRunners = append(childRunners, runner.newChild(variables))
	}
	return childRunners
}
