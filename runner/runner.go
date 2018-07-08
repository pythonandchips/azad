package runner

import (
	"github.com/pythonandchips/azad/conn"
	"github.com/zclconf/go-cty/cty"
)

type runners []*runner

func (runners runners) Close() {
	for _, runner := range runners {
		runner.Conn.Close()
	}
}

func newRunner(conn conn.Conn, address string) *runner {
	runner := runner{
		Conn:        conn,
		Address:     address,
		taskResults: map[string]taskResult{},
	}
	return &runner
}

type runner struct {
	Address     string
	Conn        conn.Conn
	Variables   variables
	taskResults map[string]taskResult
}

type taskResult map[string]string
type variables map[string]string

func (runner *runner) setResult(name string, result map[string]string) {
	runner.taskResults[name] = result
}

func (runner runner) toContext() map[string]cty.Value {
	context := map[string]cty.Value{
		"var": mapToCtyValue(runner.Variables),
	}
	for k, v := range runner.taskResults {
		context[k] = mapToCtyValue(v)
	}
	return context
}

func mapToCtyValue(data map[string]string) cty.Value {
	if len(data) == 0 {
		return cty.MapValEmpty(cty.String)
	}
	vars := map[string]cty.Value{}
	for k, v := range data {
		vars[k] = cty.StringVal(v)
	}
	return cty.MapVal(vars)
}
