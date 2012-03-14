package gpkg

import "testing"

const SOURCE_TEST_FILE = "sources"
const SOURCE_TEST_PACKAGE = "example1"
const SOURCE_TEST_FALLBACK_PACKAGE = "ptcp"

func TestNewSources(t *testing.T) {
	sources := NewSources(SOURCE_TEST_FILE)
	if sources == nil {
		t.Fatal("Failed to load sources")
	}
}
/*
func TestSourcesFind(t *testing.T) {
	sources := NewSources(SOURCE_TEST_FILE)
	if sources == nil {
		t.Fatal("Failed to load sources")
	}
	v := sources.Find(SOURCE_TEST_PACKAGE)
	if v == nil {
		t.Fatal("Failed to find package")
	}
}

func TestSourcesFindFallback(t *testing.T) {
	sources := NewSources(SOURCE_TEST_FILE)
	if sources == nil {
		t.Fatal("Failed to load sources")
	}
	v := sources.Find(SOURCE_TEST_FALLBACK_PACKAGE)
	if v == nil {
		t.Fatal("Failed to find package")
	}
}*/
