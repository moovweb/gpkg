package builder

import "os"
import "fmt"
import "path/filepath"
import "errors"

import . "github.com/moovweb/gpkg/container"
import . "github.com/moovweb/gpkg/pkg"
import . "github.com/moovweb/gpkg/specs"
import . "github.com/moovweb/gpkg/source"
import . "github.com/moovweb/gpkg/util"

type Builder struct {
	status  int
	sources *Sources
	pkg     *Package
	specs   *Specs
	deps    map[string]*Package

	build_container  Container
	import_container Container
	source_container Container
}

const (
	INITIAL   = iota
	CLONED    // Cloned source state unknown
	CLEAN     // Clean source is clean
	READY     // Dependencies loaded
	BUILT     // Built and ready for testing or install
	TESTED    // Built and tested ready for install
	INSTALLED // Installed into the local gpkg pkgset
	ERROR     // Something went wrong :(
)

func NewBuilder(sources *Sources, name string, spec string, tmpdir string) *Builder {
	b := &Builder{
		sources:          sources,
		status:           INITIAL,
		deps:             map[string]*Package{},
		build_container:  NewSimpleContainer(filepath.Join(tmpdir, "build")),
		import_container: NewSimpleContainer(filepath.Join(tmpdir, "import")),
		source_container: NewSimpleContainer(tmpdir),
	}

	f, v, source := sources.FindInSources(name, spec)
	if f == false {
		return nil
	}
	b.pkg = NewPackage(name, v, source)
	return b
}

func pushGopath(path string, preserve bool) string {
	old_gopath := os.Getenv("GOPATH")
	gopath := path
	if preserve == true && old_gopath != "" {
		gopath = gopath + ":" + old_gopath
	} else {
	}
	os.Setenv("GOPATH", gopath)
	return old_gopath
}

func (b *Builder) LoadImports() error {
	for name, spec := range b.specs.List {
		f, v, source := b.sources.FindInCache(name, spec)
		if f == true {
			b.sources.LoadFromCache(name, v, b.import_container.String())
		} else {
			//newbuilder := NewBuilder(b.sources, name, spec, source_container.String())
			return errors.New("Couldn't find package in cache")
		}
		fmt.Println(f, v, source)
	}
	return nil
}

func (b *Builder) Clone() error {
	err := b.pkg.Clone(b.source_container)
	if err == nil {
		b.status = CLONED
	} else {
		b.status = ERROR
		return err
	}
	specs, err := NewSpecs(filepath.Join(b.source_container.SrcDir(), b.pkg.Name))
	if err != nil {
		b.status = ERROR
		return err
	}
	if specs != nil {
		b.specs = specs
	} else {
		b.specs = NewBlankSpecs(b.pkg.Source)
	}
	fmt.Println(b.specs)
	return err
}

func (b *Builder) Clean() (string, error) {
	out, err := b.pkg.Clean()
	if err == nil {
		b.status = CLEAN
	} else {
		b.status = ERROR
	}
	return out, err
}

func (b *Builder) Build() (string, error) {
	err := b.LoadImports()
	if err != nil {
		return "", err
	}

	gopath := pushGopath(b.import_container.String(), false)
	out, err := b.pkg.Build()
	pushGopath(gopath, false)

	if err == nil {
		b.status = BUILT
	} else {
		b.status = ERROR
	}
	return out, err
}

func (b *Builder) Test() error {
	return nil
}

func (b *Builder) Install(dest Container) (string, error) {
	gopath := pushGopath(dest.String(), false)
	pushGopath(b.import_container.String(), true)
	out, err := b.pkg.Install()
	pushGopath(gopath, false)
	FileCopy(b.source_container.SrcDir(), dest.SrcDir())
	FileCopy(b.source_container.BinDir(), dest.BinDir())
	return out, err
}
