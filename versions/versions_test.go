package gpkg

import "testing"

func TestNewVersion(t *testing.T) {
	v := NewVersion("0.0.1")
	if v == nil {
		t.Fatal("Failed to create version", "0.0.1")
	}
	if v.String() != "0.0.1" {
		t.Fatal("Unexpected version string", v.String(), "!=", "0.0.1")
	}
}

func TestNewVersionFromMatch(t *testing.T) {
	versions := []Version{
		*NewVersion("0.0.1"), 
		*NewVersion("0.0.2"), 
		*NewVersion("0.0.3"),
		*NewVersion("0.1.1"), 
		*NewVersion("0.1.2"), 
		*NewVersion("0.1.3"),
		*NewVersion("1.0.1"), 
		*NewVersion("1.1.2"), 
		*NewVersion("1.4.44"),
	}

	v, err := NewVersionFromMatch(versions, "0.0.1")
	if err != nil {
		t.Fatal(err)
	}
	if v.String() != "0.0.1" {
		t.Fatal("Unexpected version string", v.String(), "!=", "0.0.1")
	}

	v, err = NewVersionFromMatch(versions, "~> 0.0.1")
	if err != nil {
		t.Fatal(err)
	}
	if v.String() != "0.0.3" {
		t.Fatal("Unexpected version string", v.String(), "!=", "0.0.3")
	}

	v, err = NewVersionFromMatch(versions, "~> 1.1.0")
	if err != nil {
		t.Fatal(err)
	}
	if v.String() != "1.1.2" {
		t.Fatal("Unexpected version string", v.String(), "!=", "1.1.2")
	}
}

