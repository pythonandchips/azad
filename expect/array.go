package expect

import "testing"

// EqualFatal asserts that length of 2 number are equal
// if if fails it is raises a fatal
func EqualFatal(t *testing.T, length, expected int, format ...string) {
	if length != expected {
		t.Fatalf("expected %d to equal %d", expected, length)
	}
}
