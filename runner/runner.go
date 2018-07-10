package runner

import (
	"github.com/pythonandchips/azad/conn"
	"github.com/zclconf/go-cty/cty"
)

func newRunner(conn conn.Conn, address string, variables map[string]cty.Value) *runner {
	runner := runner{
		Conn:            conn,
		Address:         address,
		GlobalVariables: variables,
		taskResults:     map[string]taskResult{},
	}
	return &runner
}

type runner struct {
	Address         string
	Conn            conn.Conn
	GlobalVariables map[string]cty.Value
	taskResults     map[string]taskResult
}

func (r *runner) newChild(variables map[string]cty.Value) *runner {
	return &runner{
		Address:         r.Address,
		Conn:            r.Conn,
		GlobalVariables: mergeVariables(r.GlobalVariables, variables),
		taskResults:     map[string]taskResult{},
	}
}

type taskResult map[string]string

func (taskResult taskResult) toCtyValue() cty.Value {
	if len(taskResult) == 0 {
		return cty.MapValEmpty(cty.String)
	}
	variables := map[string]cty.Value{}
	for key, value := range taskResult {
		variables[key] = cty.StringVal(value)
	}
	return cty.MapVal(variables)
}

func (r *runner) setResult(name string, result map[string]string) {
	r.taskResults[name] = result
}

func (r runner) toContext() map[string]cty.Value {
	context := map[string]cty.Value{
		"var": mapToCtyValue(r.GlobalVariables),
	}
	for k, v := range r.taskResults {
		context[k] = v.toCtyValue()
	}
	return context
}

func mapToCtyValue(data map[string]cty.Value) cty.Value {
	if len(data) == 0 {
		return cty.MapValEmpty(cty.String)
	}
	return cty.MapVal(data)
}
