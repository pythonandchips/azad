package stringslice

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExists(t *testing.T) {
	slice := []string{"hello", "world"}
	t.Run("returns the index of matching string", func(t *testing.T) {
		assert.Equal(t, Exists("hello", slice), true)
	})
	t.Run("returns the -1 if no match found", func(t *testing.T) {
		assert.Equal(t, Exists("goodbye", slice), false)
	})
}
