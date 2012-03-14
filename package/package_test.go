package gpkg

import "testing"
import . "sources"

func TestNewPackage(t *testing.T) {
	p := NewPackage("test-package")
	if p == nil {
		t.Fatal("Failed to create package")
	}
}

func TestNewPackageFromSource(t *testing.T) {
	s := NewSource("git@github.com:moovweb")
	p := NewPackageFromSource("test-package", s)

	if p == nil {
		t.Fatal("Failed to create package")
	}
	if p.Source == nil {
		t.Fatal("Failed to set Source")
	}
}
