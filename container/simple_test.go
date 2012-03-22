package container

import "testing"

func TestNewSimpleContainer(t *testing.T) {
	c := NewSimpleContainer("/tmp/gpkg-container-test")
	err := c.Empty()
	if err != nil {
		t.Fatal(err)
	}
}

