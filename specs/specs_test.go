package specs

import "testing"

import "io/ioutil"

const TEST_FILE = "Package.gvm"

func TestNewSpecs(t *testing.T) {
	_, err := NewSpecs(TEST_FILE + ".bogus")
	if err == nil {
		t.Fatal("Invalid filename not returning nil")
	}

	specs, err := NewSpecs(TEST_FILE)
	if err != nil {
		t.Error(err)
		t.Fatal("Failed to load test specs from " + TEST_FILE)
	}

	if specs == nil {
		t.Fatal("Specs is nil")
	}
	if specs.String() == "" {
		t.Fatal("Specs to string is blank")
	}

	out, ioerr := ioutil.ReadFile(TEST_FILE)
	if ioerr != nil {
		t.Fatal("Failed to read test spec file")
	}
	if specs.String() != string(out) {
		t.Fatal("Render back to text didn't match", "\n" + specs.String(), "\n" + string(out))
	}
}
