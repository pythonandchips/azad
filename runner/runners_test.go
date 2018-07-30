package runner

import (
	"fmt"
	"testing"

	"github.com/pythonandchips/azad/conn"
	"github.com/pythonandchips/azad/expect"
	"github.com/stretchr/testify/assert"
	"github.com/zclconf/go-cty/cty"
)

func TestCreateRunners(t *testing.T) {
	addresses := []string{"10.0.0.1", "10.0.0.2"}
	variables := map[string]cty.Value{
		"install_dir": cty.StringVal("/opt"),
	}
	t.Run("when connection can be made to address", func(t *testing.T) {
		newConn = func(config conn.Config) conn.Conn {
			return &fakeConn{}
		}
		runners, err := createRunners(addresses, variables)
		if err != nil {
			t.Fatalf("Unexpected Error: %s", err)
		}
		t.Run("creates runners for addresses", func(t *testing.T) {
			expect.EqualFatal(t, len(runners), 2)
			RunnerEquals(t, runners[0], "10.0.0.1", variables)
			RunnerEquals(t, runners[1], "10.0.0.2", variables)
		})
	})
	t.Run("when connection cannot be made to server an error is returned", func(t *testing.T) {
		newConn = func(config conn.Config) conn.Conn {
			return &fakeConn{ConnectionError: fmt.Errorf("unable to connect")}
		}
		_, err := createRunners(addresses, variables)
		if err == nil {
			t.Fatalf("Expected error but got none")
		}
	})
}

func TestRunnersChildRunners(t *testing.T) {
	runner := &runner{
		Address: "10.0.0.1",
		Conn:    &fakeConn{Closed: false},
		GlobalVariables: map[string]cty.Value{
			"install_dir": cty.StringVal("/opt"),
		},
	}
	runners := runners{runner}
	childVariables := map[string]cty.Value{
		"install_ruby": cty.StringVal("/opt/ruby"),
	}
	childRunners := runners.ChildRunners(childVariables)
	expect.EqualFatal(t, len(childRunners), 1)

	childRunner := childRunners[0]

	assert.Equal(t, childRunner.GlobalVariables["install_dir"].AsString(), "/opt")
	assert.Equal(t, childRunner.GlobalVariables["install_ruby"].AsString(), "/opt/ruby")
	assert.Equal(t, childRunner.Address, runner.Address)
	assert.Equal(t, childRunner.Conn, runner.Conn)
}

func TestRunnersClose(t *testing.T) {
	runners := runners{
		&runner{Address: "10.0.0.1", Conn: &fakeConn{Closed: false}},
		&runner{Address: "10.0.0.2", Conn: &fakeConn{Closed: false}},
	}
	runners.Close()
	t.Run("closes all connections", func(t *testing.T) {
		assert.True(t, runners[0].Conn.(*fakeConn).Closed)
		assert.True(t, runners[1].Conn.(*fakeConn).Closed)
	})
}
