package runner

import (
	"strings"

	multierror "github.com/hashicorp/go-multierror"
	"github.com/hashicorp/hcl2/hcl"
	"github.com/jinzhu/copier"
	"github.com/pythonandchips/azad/conn"
	"github.com/zclconf/go-cty/cty"
)

type contextStore struct {
	store
	contextUser        string
	contextConnections connections
	contextPath        string
}

func newContextStore(applyTo []string, user, path string, baseStore store) (contextStore, error) {
	contextStore := contextStore{contextPath: path, contextUser: user, store: baseStore}
	servers := contextStore.serversForGroups(applyTo)
	for _, server := range servers {
		conn := conn.NewConn(baseStore.config.SSHConfig())
		err := conn.ConnectTo(server.address)
		if err != nil {
			return contextStore, err
		}
		contextStore.contextConnections = append(contextStore.contextConnections, &connection{
			taskVariables: variables{},
			conn:          conn,
		})
	}
	return contextStore, nil
}

func (contextStore contextStore) evalVariableForTask(
	attr *hcl.Attribute,
	connection *connection,
	additionalVariables map[string]cty.Value,
	allowedTypes []string,
) (cty.Value, error) {
	evalContext, err := contextStore.evalContextForTask(attr.Expr.Variables(), connection, additionalVariables)
	if err != nil {
		return cty.StringVal(""), err
	}
	val, evalErr := attr.Expr.Value(evalContext)
	if evalErr != nil {
		return cty.StringVal(""), evalErr
	}
	return val, nil
}

func (contextStore contextStore) evalContextForTask(requiredVariables []hcl.Traversal, connection *connection, variables map[string]cty.Value) (*hcl.EvalContext, error) {
	errors := &multierror.Error{}
	var err error
	for _, requiredVariable := range requiredVariables {
		variablePath := variablePathFromTraverser(requiredVariable)
		switch variablePath[0] {
		case "var":
			variables, err = contextStore.variableFromPath(variables, variablePath)
			if err != nil {
				errors = multierror.Append(errors, err)
			}
		case "srv":
			variables, err = contextStore.serverFromPath(variables, variablePath)
			if err != nil {
				errors = multierror.Append(errors, err)
			}
		case "item":
			// special case do nothing any item will be initialize at the top
		default:
			variableName := strings.Join(variablePath, ".")
			variable, err := connection.findVariable(variableName)
			if err != nil {
				errors = multierror.Append(errors, err)
				continue
			}
			variables = buildVariable(variablePath, variable.value, variables)
		}
	}
	return &hcl.EvalContext{
		Variables: variables,
	}, errors.ErrorOrNil()
}

func (c contextStore) childStore(user, file string) contextStore {
	childContext := contextStore{}
	copier.Copy(&childContext, &c)
	if user != "" {
		childContext.contextUser = user
	}
	childContext.contextPath = file
	return childContext
}

func (contextStore contextStore) user() string {
	return contextStore.contextUser
}

func (contextStore contextStore) playbookPath() string {
	return contextStore.path
}

func (contextStore contextStore) rolePath() string {
	return contextStore.contextPath
}

func (contextStore contextStore) connections() connections {
	return contextStore.contextConnections
}

func (contextStore contextStore) closeConnections() {
	for _, connection := range contextStore.contextConnections {
		connection.close()
	}
}

type connections []*connection

func (connections connections) each(f func(*connection) error) error {
	errors := &multierror.Error{}
	for _, currentConnection := range connections {
		err := f(currentConnection)
		if err != nil {
			errors = multierror.Append(errors, err)
		}
	}
	return errors.ErrorOrNil()
}

type connection struct {
	taskVariables variables
	conn          conn.Conn
}

func (connection *connection) addTaskVariable(variable variable) {
	connection.taskVariables = append(connection.taskVariables, variable)
}

func (connection *connection) findVariable(name string) (variable, error) {
	return connection.taskVariables.findByName(name)
}

func (connection *connection) close() {
	connection.conn.Close()
}
