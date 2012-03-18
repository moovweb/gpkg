package pkg

import "exec"
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
	gvm *Gvm
	root string
	name string
	tag string
	version *Version
	Source
	tmpdir string
	logger *Logger
	deps map[string]*Package

	specs *Specs
	tool Tool
	force_install bool
}

func NewPackage(gvm *Gvm, name string, tag string, root string, Source Source, tmpdir string, logger *Logger) *Package {
	p := &Package{
		root: root,
		gvm: gvm,
		logger: logger,
		name: name,
		tag: tag,
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

func (p *Package) Get() os.Error {
	tmp_src_dir := filepath.Join(p.tmpdir, p.name, "src")
	os.RemoveAll(tmp_src_dir)
	os.MkdirAll(tmp_src_dir, 0775)
	p.logger.Debug(" * Getting Source")
	err := p.Clone(p.name, p.tag, tmp_src_dir)
	if err != nil {
		return err
	}

	if p.tag != "" {
		p.logger.Debug(" * Checking out ", p.tag)
		err := os.Chdir(tmp_src_dir + "/" + p.name)
		if err != nil {
			p.logger.Fatal("Unable to chdir to checkout version", p.tag, "of", p.name)
		}
		_, err = exec.Command("git", "checkout", p.tag).CombinedOutput()
		if err != nil {
			p.logger.Error("Invalid version:", p.tag, "of", p.name, "specified")
		}
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
			/*err := os.Chdir(tmp_src_dir + "/" + p.name)
			out, err := exec.Command("git", "describe", "--exact-match", "--tags", "--match", "*.*.*").CombinedOutput()
			if err == nil {
				p.tag = strings.TrimSpace(string(out))
			} else {
				if p.tag != "" {
					p.tag += ".src"
				} else {
					p.tag = "src"
				}
			}*/
			p.tag += ".src"
		}
	}

	return nil
}

func (p *Package) LoadImport(dep *Package, dir string) {
	//p.logger.Trace("dep", fmt.Sprint(dep))
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

	//p.deps += "pkg " + dep.name + " " + dep.tag + "\n"
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
		if spec != "" {
			found, source := p.gvm.FindPackageByVersion(name, spec)
			if found == true {
				dep = NewPackage(p.gvm, name, spec, source, NewSource(source), p.tmpdir, p.logger)
			}
		} else {
			found, version, source := p.gvm.FindPackage(name)
			if found == true {
				p.specs.List[name] = version
				dep = NewPackage(p.gvm, name, version, source, NewSource(source), p.tmpdir, p.logger)
			}
		}

		if dep == nil {
			p.logger.Fatal("ERROR: Couldn't find " + name + " " + spec + " in any sources")
		}
		p.logger.Debug("    -", dep.name, dep.tag, "(Spec:", spec + ")")
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
	tmp_src_dir := filepath.Join(p.tmpdir, p.name, "src")

	if !p.LoadImports(tmp_import_dir) {
		p.logger.Error("Failed to load imports")
		return false
	}

	os.Setenv("GOPATH", tmp_build_dir + ":" + tmp_import_dir)
	old_build_number := os.Getenv("BUILD_NUMBER")	
	os.Setenv("BUILD_NUMBER", p.tag)
	_, err := os.Open(filepath.Join(tmp_src_dir, p.name, "Makefile.gvm"))
	if err == nil {
		p.tool = NewMakeTool(filepath.Join(tmp_src_dir, p.name), "Makefile.gvm")
	} else {
		p.tool = NewGbTool(filepath.Join(tmp_src_dir, p.name))
	}
	out, berr := p.tool.Clean()
	if berr != nil {
		p.logger.Error("Failed to clean")
		p.logger.Error(p.PrettyLog(out))
		return false
	}
	// Build using gpkg Makefile
	p.logger.Debug(" * Building")
	out, berr = p.tool.Build()
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
	// Run tests using gpkg Makefile
	p.logger.Debug(" * Installing")
	out, berr = p.tool.Install()
	if berr != nil {
		p.logger.Error(berr)
		p.logger.Error(p.PrettyLog(out))
		return false
	} else {
		p.logger.Debug(p.PrettyLog(out))
	}

	os.Setenv("BUILD_NUMBER", old_build_number)

	p.logger.Debug(" * Installing", p.name + "-" + p.tag + "...")

	//if p.force_install == true {
		err = os.RemoveAll(filepath.Join(p.root, p.tag))
		if err != nil {
			p.logger.Fatal("Failed to remove old version")
		}
	/*} else {
		_, err := os.Open(filepath.Join(p.root, p.tag))
		if err == nil {
			p.logger.Fatal("Already installed!")
		}
	}*/
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
	return true
}

func (p *Package) Install(tmpdir string) {
	p.logger.Debug("Starting install of", p.name, p.tag)
	if p.Source == nil {
		p.logger.Fatal("No source specified")
	}
	err := p.Get()
	if err != nil {
		p.logger.Fatal("ERROR Getting package source", err)
	}
	if !p.Build() {
		// TODO: p.gvm.DeletePackage(p)
		p.logger.Fatal("ERROR Building package")
	}
}
