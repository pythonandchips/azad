package azad

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPluginName(t *testing.T) {
	t.Run("with valid inventory name", func(t *testing.T) {
		inventory := Inventory{Name: "aws.ec2"}
		name, err := inventory.PluginName()
		assert.Nil(t, err)
		assert.Equal(t, name, "aws")
	})
	t.Run("with invalid inventory name", func(t *testing.T) {
		inventory := Inventory{Name: "ec2"}
		_, err := inventory.PluginName()
		assert.NotNil(t, err)
	})
}

func TestServiceName(t *testing.T) {
	t.Run("with valid inventory name", func(t *testing.T) {
		inventory := Inventory{Name: "aws.ec2"}
		name, err := inventory.ServiceName()
		assert.Nil(t, err)
		assert.Equal(t, name, "ec2")
	})
	t.Run("with invalid inventory name", func(t *testing.T) {
		inventory := Inventory{Name: "ec2"}
		_, err := inventory.ServiceName()
		assert.NotNil(t, err)
	})
}
