package stringslice

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIndexOf(t *testing.T) {
	slice := []string{"hello", "world"}
	t.Run("returns the index of matching string", func(t *testing.T) {
		assert.Equal(t, IndexOf("hello", slice), 0)
	})
	t.Run("returns the -1 if no match found", func(t *testing.T) {
		assert.Equal(t, IndexOf("goodbye", slice), -1)
	})
}
