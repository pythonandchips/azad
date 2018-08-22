package runner

import (
	"testing"

	"github.com/pythonandchips/azad/conn"
	"github.com/pythonandchips/azad/logger"
	"github.com/pythonandchips/azad/steps"
	"github.com/stretchr/testify/assert"
	"github.com/zclconf/go-cty/cty"
)

func TestHandleIncludes(t *testing.T) {
	logger.StubLogger()
	includeStep, contextStore := setupTestHanldeIncludes()
	t.Run("returns no errors when valid", func(t *testing.T) {
		err := handleIncludes(includeStep, contextStore)
		if err != nil {
			t.Fatalf("unexpected err: %s", err)
		}
	})
	t.Run("when roles cannot be evaluated returns error", func(t *testing.T) {
		includeStep.Roles = testExpression("role", `not_a_value`)
		err := handleIncludes(includeStep, contextStore)
		errorMessage := "1 error occurred:\n\n* variable not found: not_a_value"
		assert.Equal(t, errorMessage, err.Error())
	})
	t.Run("when roles is not an array returns error", func(t *testing.T) {
		includeStep.Roles = testExpression("role", `"not_a_value"`)
		err := handleIncludes(includeStep, contextStore)
		errorMessage := `Unexpected type for role, expected [list] but was string
	Filename: , {0 0 0}-{0 0 0}`
		assert.Equal(t, err.Error(), errorMessage)
	})
	t.Run("when roles does not exist", func(t *testing.T) {
		includeStep.Roles = testExpression("role", `["not_a_value"]`)
		err := handleIncludes(includeStep, contextStore)
		errorMessage := `Container with name not_a_value not found in /home/user/azad/roles/ruby`
		assert.Equal(t, err.Error(), errorMessage)
	})
}

func setupTestHanldeIncludes() (steps.IncludesStep, *contextStore) {
	roleList = steps.RoleContainers{
		{Name: "basic_security"},
	}
	store := store{
		variables: variables{},
		servers:   servers{},
	}
	conn := &conn.LoggerSSHConn{}
	connection := connection{
		taskVariables: variables{
			{name: "install-erlang.success", value: cty.StringVal("true"), valueType: "string"},
		},
		conn: conn,
	}
	contextStore := contextStore{
		store:              store,
		contextUser:        "deploy",
		contextPath:        "/home/user/azad/roles/ruby",
		contextConnections: connections{&connection},
	}
	includeStep := steps.IncludesStep{
		Roles: testExpression("role", `["basic_security"]`),
	}
	return includeStep, &contextStore
}
