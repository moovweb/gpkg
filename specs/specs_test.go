package gpkg

import "testing"

const TEST_FILE = "Package.gvm"

func TestNewSpecs(t *testing.T) {
	_, err := NewSpecs(TEST_FILE)
	if err != nil {
		t.Error(err)
		t.Fatal("Failed to load test specs from " + TEST_FILE)
	}
}
