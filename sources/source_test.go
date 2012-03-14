package gpkg

import "testing"

func TestNewSource(t *testing.T) {
	s := NewSource("/tmp/phony")
	if s == nil {
		t.Fatal("Failed to create source")
	}
}
