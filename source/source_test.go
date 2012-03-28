package source

import "testing"

import "os"
import "path"

const TMP_TEST_ROOT = "/tmp/gpkg-test"
const TMP_TEST_DEPTH = "/tmp/gpkg-test/deep/inside/a/path"

var SourceCloneTests = map[string]string{
/*	"git://github.com/jbussdieker": "example1", 
	"https://github.com/jbussdieker": "example1", 
	"git@github.com:jbussdieker": "example1", 
	"/home/jbussdieker/moovweb/gohattan/src": "gokogiri",*/
}

func TestNewSource(t *testing.T) {
	s := NewSource("/tmp/blah")
	switch v := s.(type) {
	case LocalSource:
		break
	default:
		t.Fatal("Unexpected type for /tmp/blah")
	}
	s = NewSource("git@github.com:tmp/blah")
	switch v := s.(type) {
	case GitSource:
		break
	default:
		t.Fatal("Unexpected type for git@github.com:tmp/blah")
	}
}

func testSourceClone(source string, name string, t *testing.T) {
	s := NewSource(source)
	serr := s.Clone(name, nil, TMP_TEST_DEPTH)
	if serr != nil {
		t.Fatal("Failed to clone", name, "from", source, serr)
	}
	_, err := os.Open(path.Join(TMP_TEST_DEPTH, name))
	if err != nil {
		t.Fatal("Failed to clone", name, "from", source, err)
	}
	_, err = os.Open(path.Join(TMP_TEST_DEPTH, name, name))
	if err == nil {
		t.Fatal("Extra depth in clone", name, "from", source, err)
	}
}

func TestSourceCloneTests(t *testing.T) {
	for source, repo := range SourceCloneTests {
		err := os.RemoveAll(TMP_TEST_ROOT)
		if err != nil {
			t.Fatal("Failed to clear temp test folder", err)
		}
		testSourceClone(source, repo, t)
	}
}
