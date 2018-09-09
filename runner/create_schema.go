package runner

import (
	"github.com/hashicorp/hcl2/hcl"
	"github.com/pythonandchips/azad/plugin"
)

func createSchemaFromFields(fields []plugin.Field) *hcl.BodySchema {
	attributes := []hcl.AttributeSchema{}
	for _, field := range fields {
		attributes = append(attributes, hcl.AttributeSchema{
			Name: field.Name, Required: field.Required,
		})
	}
	return &hcl.BodySchema{
		Attributes: attributes,
	}
}
