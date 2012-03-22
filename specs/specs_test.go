package specs

import "testing"

import "os"
import "io/ioutil"
import "path/filepath"

const TEST_FILE = "Package.gvm"

func TestNewSpecs(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	specs, err := NewSpecs(wd + "/bogus")
	if err != nil || specs != nil {
		t.Fatal("Invalid filename not returning nil")
	}

	specs, err = NewSpecs(wd)
	if err != nil {
		t.Error(err)
		t.Fatal("Failed to load test specs from " + wd)
	}

	if specs == nil {
		t.Fatal("Specs is nil")
	}
	if specs.String() == "" {
		t.Fatal("Specs to string is blank")
	}

	out, ioerr := ioutil.ReadFile(filepath.Join(wd, TEST_FILE))
	if ioerr != nil {
		t.Fatal("Failed to read test spec file")
	}
	if specs.String() != string(out) {
		t.Fatal("Render back to text didn't match", "\n"+specs.String(), "\n"+string(out))
	}
}
