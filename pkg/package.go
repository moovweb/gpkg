package pkg

import "os"
import "io/ioutil"
import "path/filepath"
import "strings"

import . "gvm"
import . "logger"
import . "source"
import . "version"
import . "specs"
import . "tools"
import . "util"

type Package struct {
	gvm     *Gvm
	root    string
	name    string
	tag     string
	version *Version
	Source
	tmpdir string
	logger *Logger
	deps   map[string]*Package

	specs *Specs
	tool  Tool
}

func NewPackage(gvm *Gvm, name string, tag string, root string, Source Source, tmpdir string, logger *Logger) *Package {
	p := &Package{
		root:   root,
		gvm:    gvm,
		logger: logger,
		name:   name,
		tag:    tag,
		Source: Source,
		tmpdir: tmpdir,
	}
	return p
}

func (p *Package) String() string {
	return "   root: " + p.root + "\n" +
		"   name: " + p.name + "\n" +
		"    tag: " + p.tag + "\n" +
		" source: " + p.Source.Root() + "\n" +
		" tmpdir: " + p.tmpdir + "\n"
}

func (p *Package) Clone() os.Error {
	tmp_src_dir := filepath.Join(p.tmpdir, p.name, "src")
	err := p.Source.Clone(p.name, p.tag, tmp_src_dir)
	if err != nil {
		return err
	}

	if p.tag == "" {
		v, err := ioutil.ReadFile(filepath.Join(tmp_src_dir, p.name, "VERSION"))
		if err == nil {
			p.tag = strings.TrimSpace(string(v))
		} else {
			p.tag = "0.0"
		}
		if os.Getenv("BUILD_NUMBER") != "" {
			p.tag += "." + os.Getenv("BUILD_NUMBER")
		} else {
			p.tag += ".src"
		}
	}

	return nil
}

func (p *Package) LoadImport(dep *Package, dir string) {
	if p.name == dep.name {
		p.logger.Fatal("Packages cannot import themselves")
	}
	if p.deps[dep.name] != nil {
		if p.deps[dep.name].tag != dep.tag {
			p.logger.Error("Version conflict!")
			p.logger.Fatal(dep.name, "version:", p.deps[dep.name].tag, "Imported already for", p.name)
		}
	} else {
		for _, subdep := range dep.deps {
			p.LoadImport(subdep, dir)
		}
		p.deps[dep.name] = dep
	}

	os.MkdirAll(dir, 0775)
	err := FileCopy(filepath.Join(dep.root, "pkg"), dir)
	if err != nil {
		p.logger.Fatal("ERROR: Couldn't load import: " + dep.name + "\n" + err.String())
	}
}

func (p *Package) LoadImports(dir string) bool {
	tmp_src_dir := filepath.Join(p.tmpdir, p.name, "src")
	specs, err := NewSpecs(filepath.Join(tmp_src_dir, p.name, "Package.gvm"))
	if err != nil {
		p.logger.Debug(" * No dependencies found")
		p.specs = NewBlankSpecs(p.Source.Root())
		return true
	} else {
		p.specs = specs
	}

	p.deps = map[string]*Package{}

	p.logger.Debug(" * Loading dependencies for", p.name)
	for name, spec := range p.specs.List {
		var dep *Package
		found, version, source := p.gvm.FindPackage(name, spec)
		if found == true {
			dep = NewPackage(p.gvm, name, version, source, NewSource(source), p.tmpdir, p.logger)
		} else {
			found, versions, source := p.gvm.FindSource(name, spec)
			if found == true {
				v, err := NewVersionFromMatch(versions, spec)
				if err != nil {
					return false
				}
				dep = NewPackage(p.gvm, name, v.String(), filepath.Join(p.gvm.PkgsetRoot(), "pkg.gvm", name), source, p.tmpdir, p.logger)
				p.logger.Debug(dep)
				dep.Install("")
				dep.root = filepath.Join(p.gvm.PkgsetRoot(), "pkg.gvm", name, v.String())
			}
		}

		if dep == nil {
			p.logger.Fatal("ERROR: Couldn't find " + name + " " + spec + " in any sources")
		}
		p.logger.Debug("    -", dep.name, dep.tag, "(Spec:", spec+")")
		p.LoadImport(dep, dir)
	}
	return true
}

func (p *Package) WriteManifest() {
	p.specs.Origin = p.Source.Root()
	manifest := p.specs.String()
	err := ioutil.WriteFile(filepath.Join(p.root, p.tag, "manifest"), []byte(manifest), 0664)
	if err != nil {
		p.logger.Fatal("Failed to write manifest file")
	}
}

func (p *Package) PrettyLog(buf string) (formatted string) {
	lines := strings.Split(buf, "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		if formatted != "" {
			formatted += "\n"
		}
		formatted += "    : " + line
	}
	return
}

func (p *Package) Build() bool {
	tmp_build_dir := filepath.Join(p.tmpdir, p.name, "build")
	tmp_import_dir := filepath.Join(p.tmpdir, p.name, "import")

	if !p.LoadImports(tmp_import_dir) {
		p.logger.Error("Failed to load imports")
		return false
	}

	p.logger.Debug(" * Building")

	os.Setenv("GOPATH", tmp_build_dir+":"+tmp_import_dir)
	old_build_number := os.Getenv("BUILD_NUMBER")
	os.Setenv("BUILD_NUMBER", p.tag)

	out, berr := p.tool.Clean()
	if berr != nil {
		p.logger.Error("Failed to clean")
		p.logger.Error(p.PrettyLog(out))
		return false
	}
	// Build using gpkg Makefile
	out, berr = p.tool.Build()
	if berr != nil {
		p.logger.Error(berr)
		p.logger.Error(p.PrettyLog(out))
		return false
	} else {
		p.logger.Debug(p.PrettyLog(out))
	}
	// Run tests using gpkg Makefile
	out, berr = p.tool.Install()
	if berr != nil {
		p.logger.Error(berr)
		p.logger.Error(p.PrettyLog(out))
		return false
	} else {
		p.logger.Debug(p.PrettyLog(out))
	}
	// Run tests using gpkg Makefile
	p.logger.Debug(" * Testing")
	out, berr = p.tool.Test()
	if berr != nil {
		p.logger.Error(berr)
		p.logger.Error(p.PrettyLog(out))
		return false
	} else {
		p.logger.Debug(p.PrettyLog(out))
	}

	os.Setenv("BUILD_NUMBER", old_build_number)
	return true
}

func (p *Package) Install(tmpdir string) {
	tmp_build_dir := filepath.Join(p.tmpdir, p.name, "build")
	tmp_src_dir := filepath.Join(p.tmpdir, p.name, "src")

	err := p.Clone()
	if err != nil {
		p.logger.Fatal("ERROR Getting package source", err)
	}
	p.tool = NewTool(filepath.Join(tmp_src_dir, p.name))
	if !p.Build() {
		// TODO: p.gvm.DeletePackage(p)
		p.logger.Fatal("ERROR Building package")
	}

	p.logger.Debug(" * Installing", p.name+"-"+p.tag+"...")

	err = os.RemoveAll(filepath.Join(p.root, p.tag))
	if err != nil {
		p.logger.Fatal("Failed to remove old version")
	}

	os.MkdirAll(filepath.Join(p.root, p.tag), 0775)

	p.WriteManifest()

	err = FileCopy(tmp_src_dir, filepath.Join(p.root, p.tag, "src"))
	if err != nil {
		p.logger.Fatal("Failed to copy source to install folder")
	}

	err = FileCopy(filepath.Join(tmp_build_dir, "pkg"), filepath.Join(p.root, p.tag, "pkg"))
	if err != nil {
		p.logger.Fatal("Failed to copy libraries to install folder")
	}

	err = FileCopy(filepath.Join(tmp_build_dir, "bin"), filepath.Join(p.root, p.tag, "bin"))
	// TODO: err = FileCopy(filepath.Join(tmp_build_dir, "bin"), filepath.Join(p.gvm.pkgset_root))
	if err == nil {
		p.logger.Debug(" * Installed binaries")
	}

	p.logger.Info("Installed", p.name, p.tag)
}
