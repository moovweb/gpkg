package gpkg

import "testing"
import "os"

func TestContainer(t *testing.T) {
	c := NewContainer()
	_, err := os.Open(c.Root())
	if err != nil {
		t.Error("NewContainer() failed to create temp directory")
	}
	c.Close()
	_, err = os.Open(c.Root())
	if err == nil {
		t.Error("NewContainer().Close() failed to delete temp directory")
	}
}

func TestContainerSrcDir(t *testing.T) {
	c := NewContainer()
	_, err := os.Open(c.SrcDir())
	if err != nil {
		t.Error("NewContainer().SrcDir() failed to create src directory")
	}
	c.Close()
}

