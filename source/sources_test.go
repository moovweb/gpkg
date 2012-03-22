package source

import "testing"
import "fmt"

func TestNewSources(t *testing.T) {
	s := NewSources("/tmp/gpkg-builder-test/pkg.gvm", "")
	s.Add("git@github.com:moovweb")
	fmt.Println(s)
	f, v, source := s.FindInCache("log4go", "~> 0.0.0")
	fmt.Println(f, v, source)
	f, v, source = s.FindInSources("ptcp", "~> 0.0.0")
	fmt.Println(f, v, source)
}
