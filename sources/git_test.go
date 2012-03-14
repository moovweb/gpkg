package gpkg

import "testing"
import "path/filepath"
import "io/ioutil"
import "strings"
import . "util"

//const TEST_SOURCE = "git@github.com:jbussdieker"
const TEST_SOURCE = "/home/jbussdieker"
const TEST_PACKAGE = "example1"
const TEST_VERSION_FILE = "VERSION"

func TestGitSourceVersions(t *testing.T) {
	gs := NewGitSource(TEST_SOURCE)
	versions, err := gs.Versions(TEST_PACKAGE)
	if err != nil {
		t.Error(err)
		t.Fatal("Failed to get version list")
	}
	if len(versions) < 1 {
		t.Fatal("No versions returned")
	}
}

func TestGitSourceClone(t *testing.T) {
	c := NewContainer()
	defer c.Close()
	gs := NewGitSource(TEST_SOURCE)
	err := gs.Clone(TEST_PACKAGE, c.SrcDir())
	if err != nil {
		t.Fatal("Failed to clone", err)
	}
	c.Close()
}

func TestGitSourceSetVersions(t *testing.T) {
	c := NewContainer()
	defer c.Close()
	gs := NewGitSource(TEST_SOURCE)
	versions, err := gs.Versions(TEST_PACKAGE)
	if err != nil {
		t.Error(err)
		t.Fatal("Failed to get version list")
	}
	if len(versions) < 1 {
		t.Fatal("No versions returned")
	}
	err = gs.Clone(TEST_PACKAGE, c.SrcDir())
	if err != nil {
		t.Error(err)
		t.Fatal("Failed to clone")
	}
	for _, version := range versions {
		gs.SetVersion(TEST_PACKAGE, c.SrcDir(), version)
		data, err := ioutil.ReadFile(filepath.Join(c.SrcDir(), TEST_PACKAGE, TEST_VERSION_FILE))
		if err != nil {
			t.Error("Error reading version file (" + TEST_VERSION_FILE + ")")
			t.Fatal(err)
		}
		test_version := strings.TrimSpace(string(data))
		if test_version != version.String() {
			t.Fatal("Didn't get expected version:", version.String(), "got", test_version, "instead")
		}
	}
	c.Close()
}

