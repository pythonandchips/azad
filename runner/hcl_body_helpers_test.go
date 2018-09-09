package runner

import "github.com/hashicorp/hcl2/hcl"

type TestBody struct {
	attributes map[string]*hcl.Attribute
	schema     *hcl.BodySchema
	err        hcl.Diagnostics
}

func (testBody *TestBody) Content(schema *hcl.BodySchema) (*hcl.BodyContent, hcl.Diagnostics) {
	testBody.schema = schema
	return &hcl.BodyContent{
		Attributes: testBody.attributes,
	}, testBody.err
}

func (testBody *TestBody) PartialContent(schema *hcl.BodySchema) (*hcl.BodyContent, hcl.Body, hcl.Diagnostics) {
	return nil, nil, nil
}

func (testBody *TestBody) JustAttributes() (hcl.Attributes, hcl.Diagnostics) {
	return testBody.attributes, nil
}

func (testBody *TestBody) MissingItemRange() hcl.Range {
	return hcl.Range{}
}
