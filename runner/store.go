package runner

import (
	"errors"
	"fmt"
	"strings"

	multierror "github.com/hashicorp/go-multierror"
	"github.com/hashicorp/hcl2/hcl"
	"github.com/pythonandchips/azad/helpers/stringslice"
	"github.com/zclconf/go-cty/cty"
)

var allowList = []string{"list"}
var allowString = []string{"string"}
var allowAll = []string{"list", "string", "map"}
var ctyStringType = cty.String
var ctyListType = cty.Tuple([]cty.Type{cty.String})
var ctyMapType = cty.Object(map[string]cty.Type{"value": cty.String})

type evalStore interface {
	evalContextFor([]hcl.Traversal) (*hcl.EvalContext, error)
	evalVariable(*hcl.Attribute, []string) (cty.Value, error)
}

type variableStore interface {
	evalStore
	addVariable(variable) error
}

type serverStore interface {
	evalStore
	addServers(servers) error
	serversForGroups([]string) servers
}

type taskStore interface {
	evalStore
	evalVariableForTask(*hcl.Attribute, *connection, []string) (cty.Value, error)
	connections() connections
	user() string
	playbookPath() string
	rolePath() string
}

type inputStore interface {
	evalStore
	addVariable(variable) error
	playbookPath() string
	rolePath() string
}

type store struct {
	servers   servers
	variables variables
	config    Config
	path      string
}

func (store store) contextStore(applyTo []string, user, path string) (contextStore, error) {
	return newContextStore(applyTo, user, path, store)
}

func (store store) playbookPath() string {
	return store.path
}

func (store store) rolePath() string {
	return store.path
}

func (store *store) addServers(servers servers) error {
	for _, server := range servers {
		existingServer, index := store.servers.findByAddress(server.address)
		if index == -1 {
			store.servers = append(store.servers, server)
			continue
		}
		for _, groupName := range server.group {
			existingServer.addGroup(groupName)
		}
		store.servers[index] = existingServer
	}
	return nil
}

func (store *store) addVariable(variable variable) error {
	if store.variables.exists(variable.name) {
		return fmt.Errorf("Variable %s already exists in store", variable.name)
	}
	store.variables = append(store.variables, variable)
	return nil
}

func (store store) evalVariable(attr *hcl.Attribute, allowedTypes []string) (cty.Value, error) {
	evalContext, err := store.evalContextFor(attr.Expr.Variables())
	if err != nil {
		return cty.StringVal(""), err
	}
	val, valErr := attr.Expr.Value(evalContext)
	if valErr != nil {
		return cty.StringVal(""), valErr
	}
	valType, err := readableCtyType(val.Type())
	if err != nil {
		return cty.StringVal(""), err
	}
	if !stringslice.Exists(valType, allowedTypes) {
		return val, fmt.Errorf(`Unexpected type for %s, expected %s but was %s
	Filename: %s, %d-%d`, attr.Name, allowedTypes, valType, attr.Range.Filename, attr.Range.Start, attr.Range.End)
	}
	return val, nil
}

func readableCtyType(t cty.Type) (string, error) {
	if t.Equals(ctyStringType) {
		return "string", nil
	}
	if t.IsTupleType() {
		return "list", nil
	}
	if t.IsObjectType() {
		return "map", nil
	}
	return "", fmt.Errorf("unexpected type: %s", t.GoString())
}

func (store store) evalContextFor(requiredVariables []hcl.Traversal) (*hcl.EvalContext, error) {
	variables := map[string]cty.Value{}
	errors := &multierror.Error{}
	var err error
	for _, requiredVariable := range requiredVariables {
		variablePath := variablePathFromTraverser(requiredVariable)
		switch variablePath[0] {
		case "var":
			variables, err = store.variableFromPath(variables, variablePath)
			if err != nil {
				errors = multierror.Append(errors, err)
			}
		case "srv":
			variables, err = store.serverFromPath(variables, variablePath)
			if err != nil {
				errors = multierror.Append(errors, err)
			}
		default:
			errors = multierror.Append(
				errors,
				fmt.Errorf("variable not found: %s", strings.Join(variablePath, ".")),
			)
		}
	}
	return &hcl.EvalContext{
		Variables: variables,
	}, errors.ErrorOrNil()
}

func (store store) variableFromPath(variables map[string]cty.Value, variablePath []string) (map[string]cty.Value, error) {
	variableName := strings.Join(variablePath[1:], ".")
	variable, err := store.variables.findByName(variableName)
	if err != nil {
		return variables, fmt.Errorf("variable not found %s", variableName)
	}
	variables = buildVariable(variablePath, variable.value, variables)
	return variables, nil
}

func (store store) serverFromPath(variables map[string]cty.Value, variablePath []string) (map[string]cty.Value, error) {
	servers := store.serversForGroups([]string{variablePath[1]})
	addresses := []cty.Value{}
	for _, server := range servers {
		addresses = append(addresses, cty.StringVal(server.address))
	}
	if len(addresses) == 0 {
		return buildVariable(variablePath, cty.ListValEmpty(cty.String), variables), nil
	}
	return buildVariable(variablePath, cty.ListVal(addresses), variables), nil
}

func buildVariable(path []string, value cty.Value, variableMap map[string]cty.Value) map[string]cty.Value {
	if len(path) == 1 {
		variableMap[path[0]] = value
		return variableMap
	}
	if v, ok := variableMap[path[0]]; ok {
		existingVariableMap := v.AsValueMap()
		variableMap[path[0]] = mapToCtyValue(
			buildVariable(path[1:], value, existingVariableMap),
		)
		return variableMap
	}
	variableMap[path[0]] = mapToCtyValue(
		buildVariable(path[1:], value, map[string]cty.Value{}),
	)
	return variableMap
}

func mapToCtyValue(data map[string]cty.Value) cty.Value {
	return cty.ObjectVal(data)
}

func variablePathFromTraverser(traverser []hcl.Traverser) []string {
	variableParts := []string{}
	for _, part := range traverser {
		switch v := part.(type) {
		case hcl.TraverseRoot:
			variableParts = append(variableParts, v.Name)
		case hcl.TraverseAttr:
			variableParts = append(variableParts, v.Name)
		}
	}
	return variableParts
}

func (store store) serversForGroups(groupNames []string) servers {
	servers := servers{}
	for _, server := range store.servers {
		for _, group := range server.group {
			if stringslice.Exists(group, groupNames) {
				if !servers.exists(server.address) {
					servers = append(servers, server)
				}
			}
		}
	}
	return servers
}

type servers []server

var errServerNotFound = errors.New("server not found")

func (servers servers) exists(address string) bool {
	_, index := servers.findByAddress(address)
	return index != -1
}

func (servers servers) findByAddress(address string) (server, int) {
	for index, server := range servers {
		if server.address == address {
			return server, index
		}
	}
	return server{}, -1
}

type server struct {
	address string
	group   []string
}

func (server *server) addGroup(name string) {
	if server.groupExists(name) {
		return
	}
	server.group = append(server.group, name)
}

func (server server) groupExists(name string) bool {
	for _, group := range server.group {
		if group == name {
			return true
		}
	}
	return false
}

type variables []variable

var errVariableNotFound = errors.New("variable not found")

func (variables variables) findByName(name string) (variable, error) {
	for _, variable := range variables {
		if variable.name == name {
			return variable, nil
		}
	}
	return variable{}, errVariableNotFound
}

func (variables variables) exists(name string) bool {
	for _, variable := range variables {
		if variable.name == name {
			return true
		}
	}
	return false
}

type variable struct {
	name      string
	value     cty.Value
	valueType string
}
