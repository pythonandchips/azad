package runner

import (
	"testing"

	"github.com/pythonandchips/azad/expect"
	"github.com/pythonandchips/azad/logger"
	"github.com/pythonandchips/azad/steps"
	"github.com/stretchr/testify/assert"
)

func TestHandleServer(t *testing.T) {
	logger.StubLogger()
	store := store{
		variables: variables{},
		servers:   servers{},
	}
	t.Run("add servers to store", func(t *testing.T) {
		serverStep := steps.ServerStep{
			Name: "kibana_server",
			Addresses: testExpression(
				"addresses",
				`["10.0.0.1", "10.0.0.2"]`,
			),
		}
		err := handleServer(serverStep, &store)
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		expect.EqualFatal(t, len(store.servers), 2)
		assert.Equal(t, store.servers[0].address, "10.0.0.1")
		assert.Equal(t, store.servers[0].group, []string{"kibana_server"})
		assert.Equal(t, store.servers[1].address, "10.0.0.2")
		assert.Equal(t, store.servers[1].group, []string{"kibana_server"})
	})
	t.Run("returns an error if addresses cannot be evaluated", func(t *testing.T) {
		serverStep := steps.ServerStep{
			Name: "kibana_server",
			Addresses: testExpression(
				"addresses",
				`"10.0.0.1"`,
			),
		}
		err := handleServer(serverStep, &store)
		expectedMessage := `Unexpected type for addresses, expected [list] but was string
	Filename: , {0 0 0}-{0 0 0}`
		assert.Equal(t, err.Error(), expectedMessage)
	})
}
