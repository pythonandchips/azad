package expect

import "testing"

// NoErrors fails with fatal if error is present
func NoErrors(t *testing.T, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
}
