package schema

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildField(t *testing.T) {
	schema := BuildSchema()
	assert.Equal(t, schema, Schema{})
}
