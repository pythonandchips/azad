package runner

import (
	"fmt"
	"strings"

	"github.com/hashicorp/hcl2/hcl"
	"github.com/hashicorp/hcl2/hcl/hclsyntax"
	"github.com/pythonandchips/azad/steps"
)

func defaultServerStep() steps.ServerStep {
	return serverStep("kibana_server", []string{"10.0.0.1"})
}

func serverStep(name string, addresses []string) steps.ServerStep {
	addrs := []string{}
	for _, addr := range addresses {
		addrs = append(addrs, fmt.Sprintf(`"%s"`, addr))
	}
	addrExpr := fmt.Sprintf(`[%s]`, strings.Join(addrs, ","))
	return steps.ServerStep{
		Name: name,
		Addresses: testExpression(
			"addresses", addrExpr,
		),
	}
}

func defaultInventoryStep() steps.InventoryStep {
	return steps.InventoryStep{
		Type: "aws.ec2",
		Body: &TestBody{
			attributes: map[string]*hcl.Attribute{
				"access_id":  testExpression("access_id", `"abcdef123456"`),
				"secret_key": testExpression("secret_key", `"ABCDEF09876543"`),
			},
		},
	}
}

func defaultVariableStep() steps.VariableStep {
	return steps.VariableStep{
		Name:  "opts_path",
		Type:  testExpression("type", `"string"`),
		Value: testExpression("value", `"/opt"`),
	}
}

func mapVariableStep() steps.VariableStep {
	return steps.VariableStep{
		Name: "owners",
		Type: testExpression("type", `"map"`),
		Value: testExpression("value", `{
	user = "kyle"
	group = "ladder"
}`),
	}
}

func arrayVariableStep() steps.VariableStep {
	return steps.VariableStep{
		Name:  "groups",
		Type:  testExpression("type", `"array"`),
		Value: testExpression("value", `["kyle", "ladder"]`),
	}
}

func defaultTaskStep() steps.TaskStep {
	return steps.TaskStep{
		Type:  "apt-get.install",
		Label: "install-erlang",
		User:  testExpression("user", `"root"`),
		Body: &TestBody{
			attributes: map[string]*hcl.Attribute{
				"package": testExpression("package", `"erlang-full"`),
			},
		},
	}
}

func defaultContext() steps.ContextContainer {
	return steps.ContextContainer{
		Name:    "kibana server",
		User:    testExpression("user", `"root"`),
		ApplyTo: testExpression("apply-to", `["development", "kibana_server"]`),
		Steps: steps.StepList{
			defaultTaskStep(),
		},
	}
}

func defaultIncludeStep() steps.IncludesStep {
	return steps.IncludesStep{
		Roles: testExpression("role", `["basic_security"]`),
	}
}

func defaultRoleContainer() steps.RoleContainer {
	return steps.RoleContainer{
		Name: "basic_security",
		File: "/home/user/azad/roles/basic_security",
		User: testExpression("user", `"admin"`),
		Steps: steps.StepList{
			defaultVariableStep(),
			defaultTaskStep(),
		},
	}
}

func defaultInputTask() steps.InputStep {
	return steps.InputStep{
		Type: "core.ini",
		Body: &TestBody{
			attributes: map[string]*hcl.Attribute{
				"path": testExpression("package", `"ini_file"`),
			},
		},
	}
}

func testExpression(name, value string) *hcl.Attribute {
	expr, parseErr := hclsyntax.ParseExpression([]byte(value), "testfile.az", hcl.Pos{Line: 1, Column: 1})
	if parseErr.HasErrors() {
		panic("unable to parse test string: " + parseErr.Error())
	}
	return &hcl.Attribute{
		Name: name, Expr: expr,
	}
}
