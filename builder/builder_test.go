package builder

import "testing"
import "fmt"

import . "container"
import . "source"

func testBuilder(name string, version string, t *testing.T) *Builder {
	s := NewSources("/tmp/gpkg-builder-test/pkg.gvm", "")
	s.Add("git@github.com:moovweb")
	b := NewBuilder(s, name, version, "/tmp/gpkg-builder-test")
	return b
}

func testInstall(b *Builder, c Container, t *testing.T) {
	if b == nil {
		t.Fail()
	}
	err := b.Clone()
	if err != nil {
		t.Fatal(err)
	}
	out, err := b.Clean()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(out)
	out, err = b.Build()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(out)
	out, err = b.Install(c)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(out)
}

func TestNewBuilder(t *testing.T) {
	b := testBuilder("log4go", "0.0.94", t)
	if b == nil {
		t.Fail()
	}
}

func TestSimple(t *testing.T) {
	b := testBuilder("log4go", "0.0.94", t)
	testInstall(b, NewSimpleContainer("/tmp/gpkg-builder-test/pkg.gvm/log4go/0.0.94"), t)
}

func TestAdvanced(t *testing.T) {
	b := testBuilder("nosaka", "0.0.204", t)
	testInstall(b, NewSimpleContainer("/tmp/gpkg-builder-test/pkg.gvm/nosaka/0.0.204"), t)
}
