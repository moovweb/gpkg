package gpkg

//import "path/filepath"
//import "fmt"
import . "specs"
import . "sources"
import . "versions"
import . "tools"
import . "util"

type Package struct {
	name string

	Source
	Container
	Specs
	Tool
	Version
}

func NewPackage(name string) *Package {
	return &Package{name:name}
}

func NewPackageFromSource(name string, source Source) *Package {
	p := NewPackage(name)
	p.Source = source
	return p
}

/*
func NewPackage(name string, source Source, version *Version) *Package {
	return &Package{name: name, Source: source, version: version}
}

func (p *Package) Version() string {
	return p.version.String()
}

func (p *Package) Clone(dest Container) {
	p.Container = dest
	p.Source.Clone(p.Name(), dest.SrcDir())
	p.Specs, _ = NewSpecs(filepath.Join(p.SrcDir(), "Package.gvm"))
	//p.Tool = NewTool(p.SrcDir())
}

func (p *Package) Clean() (string, Error) {
	return p.Tool.Clean()
}

func (p *Package) Build() (string, Error) {
	return p.Tool.Build()
}

func (p *Package) Test() (string, Error) {
	return p.Tool.Test()
}

func (p *Package) SrcDir() string {
	return filepath.Join(p.Container.SrcDir(), p.name)
}

func (p *Package) Name() string {
	return p.name
}
*/
