package pkg

import "os"
import "io/ioutil"
import "path/filepath"
import "strings"
import "fmt"

import . "github.com/moovweb/gpkg/gvm"
import . "github.com/moovweb/gpkg/logger"
import . "github.com/moovweb/gpkg/source"
import . "github.com/moovweb/gpkg/version"
import . "github.com/moovweb/gpkg/specs"
import . "github.com/moovweb/gpkg/tools"
import . "github.com/moovweb/gpkg/util"

type BuildOpts struct {
	Build       bool
	Test        bool
	Install     bool
	InstallDeps bool
	UseSystem	bool
}

type Package struct {
	Source
	BuildOpts
	gvm     *Gvm
	root    string
	name    string
	version *Version
	tmpdir  string
	logger  *Logger
	deps    map[string]*Package
	specs   *Specs
	tool    Tool
}

func NewPackage(gvm *Gvm, name string, version *Version, root string, Source Source, tmpdir string, logger *Logger) *Package {
	p := &Package{
		root:    root,
		gvm:     gvm,
		logger:  logger,
		name:    name,
		version: version,
		Source:  Source,
		tmpdir:  tmpdir,
	}
	return p
}

func (p *Package) String() string {
	return fmt.Sprintln("   root:", p.root) +
		fmt.Sprintln("   name:", p.name) +
		fmt.Sprintln("version:", p.version) +
		fmt.Sprintln(" source:", p.Source) +
		fmt.Sprintln(" tmpdir:", p.tmpdir)
}

func (p *Package) Clone() os.Error {
	p.logger.Info("Downloading", p.name, p.version)
	tmp_src_dir := filepath.Join(p.tmpdir, p.name, "src")
	err := p.Source.Clone(p.name, p.version, tmp_src_dir)
	if err != nil {
		return err
	}

	if p.version == nil {
		version_str := ""
		v, err := ioutil.ReadFile(filepath.Join(tmp_src_dir, p.name, "VERSION"))
		if err == nil {
			version_str = strings.TrimSpace(string(v))
		} else {
			version_str = "0.0"
		}
		if os.Getenv("BUILD_NUMBER") != "" {
			version_str += "." + os.Getenv("BUILD_NUMBER")
		} else {
			version_str += ".src"
		}
		p.version = NewVersion(version_str)
	}

	return nil
}

func (p *Package) LoadImport(dep *Package, dir string) {
	if p.name == dep.name {
		p.logger.Fatal("Packages cannot import themselves")
	}
	if p.deps[dep.name] != nil {
		if p.deps[dep.name].version.String() != dep.version.String() {
			p.logger.Error("Version conflict!")
			p.logger.Fatal(dep.name, "version:", p.deps[dep.name].version, "Imported already for", p.name)
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
		p.logger.Fatal("ERROR: Couldn't load import lib: " + dep.name + "\n" + err.String())
	}
	err = FileCopy(filepath.Join(dep.root, "src", dep.name), filepath.Join(dir, "..", "src"))
	if err != nil {
		p.logger.Fatal("ERROR: Couldn't load import src: " + dep.name + "\n" + err.String())
	}
}

func (p *Package) LoadImports(dir string) bool {
	tmp_src_dir := filepath.Join(p.tmpdir, p.name, "src")
	specs, err := NewSpecs(filepath.Join(tmp_src_dir, p.name, "Package.gvm"))
	if err != nil {
		p.logger.Debug(" * No dependencies found")
		p.specs = NewBlankSpecs(p.Source)
		return true
	} else {
		p.specs = specs
	}

	p.deps = map[string]*Package{}

	p.logger.Info(" * Loading dependencies for", p.name)
	for name, spec := range p.specs.List {
		var dep *Package
		found, version, source := p.gvm.FindPackageInCache(name, spec)
		if found == true {
			dep = NewPackage(p.gvm, name, version, filepath.Join(p.gvm.PkgsetRoot(), "pkg.gvm", name, version.String()), source, p.tmpdir, p.logger)
		} else if p.BuildOpts.Install == true {
			found, version, source := p.gvm.FindPackageInSources(name, spec)
			if found == true {
				dep = NewPackage(p.gvm, name, version, filepath.Join(p.gvm.PkgsetRoot(), "pkg.gvm", name), source, p.tmpdir, p.logger)
				p.logger.Trace(dep)
				dep.Install(p.BuildOpts)
				dep.root = filepath.Join(p.gvm.PkgsetRoot(), "pkg.gvm", name, version.String())
			}
		} else {
			p.logger.Fatal("ERROR: Package", name , spec, "not found locally and install=false")
		}

		if dep == nil {
			p.logger.Fatal("ERROR: Couldn't find " + name + " " + spec + " in any sources")
		}
		p.logger.Info("    -", dep.name, dep.version, "(Spec:", spec+")")
		p.LoadImport(dep, dir)
	}
	return true
}

func (p *Package) WriteManifest() {
	p.specs.Origin = p.Source
	manifest := p.specs.String()
	err := ioutil.WriteFile(filepath.Join(p.root, p.version.String(), "manifest"), []byte(manifest), 0664)
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

	if p.BuildOpts.Build == true {
		p.logger.Info(" * Building")

		old_gopath := os.Getenv("GOPATH")
		gopath := tmp_build_dir+":"+tmp_import_dir
		if p.BuildOpts.UseSystem && old_gopath != "" {
			gopath = gopath+":"+old_gopath
		} else {
		}
		os.Setenv("GOPATH", gopath)

		old_build_number := os.Getenv("BUILD_NUMBER")
		os.Setenv("BUILD_NUMBER", p.version.String())

		// Clean
		out, berr := p.tool.Clean()
		if berr != nil {
			p.logger.Error("Failed to clean")
			p.logger.Error(p.PrettyLog(out))
			return false
		}
		// Build
		out, berr = p.tool.Build()
		if berr != nil {
			p.logger.Error(berr)
			p.logger.Error(p.PrettyLog(out))
			return false
		} else {
			p.logger.Info(p.PrettyLog(out))
		}
		// Install
		out, berr = p.tool.Install()
		if berr != nil {
			p.logger.Error(berr)
			p.logger.Error(p.PrettyLog(out))
			return false
		} else {
			p.logger.Info(p.PrettyLog(out))
		}
		if p.BuildOpts.Test == true {
			// Test
			p.logger.Info(" * Testing")
			out, berr = p.tool.Test()
			if berr != nil {
				p.logger.Error(berr)
				p.logger.Error(p.PrettyLog(out))
				return false
			} else {
				p.logger.Info(p.PrettyLog(out))
			}
		}

		os.Setenv("GOPATH", old_gopath)
		os.Setenv("BUILD_NUMBER", old_build_number)
	}
	return true
}

func (p *Package) Install(b BuildOpts) {
	p.BuildOpts = b

	tmp_build_dir := filepath.Join(p.tmpdir, p.name, "build")
	tmp_src_dir := filepath.Join(p.tmpdir, p.name, "src")

	err := p.Clone()
	if err != nil {
		p.logger.Fatal("ERROR Getting package source", err)
	}
	p.tool = NewTool(filepath.Join(tmp_src_dir, p.name))
	if !p.Build() {
		p.logger.Fatal("ERROR Building package")
	}

	// INSTALL
	//////////////////////
	if p.BuildOpts.Install == true {
		p.logger.Info(" * Installing", p.name+"-"+p.version.String()+"...")
		err = os.RemoveAll(filepath.Join(p.root, p.version.String()))
		if err != nil {
			p.logger.Fatal("Failed to remove old version")
		}
		os.MkdirAll(filepath.Join(p.root, p.version.String()), 0775)
		p.WriteManifest()
		err = FileCopy(tmp_src_dir, filepath.Join(p.root, p.version.String(), "src"))
		if err != nil {
			p.logger.Fatal("Failed to copy source to install folder", err)
		}

		err = FileCopy(filepath.Join(tmp_build_dir, "pkg"), filepath.Join(p.root, p.version.String(), "pkg"))
		if err != nil {
			p.logger.Fatal("Failed to copy libraries to install folder\n", err.String())
		}

		err = FileCopy(filepath.Join(tmp_build_dir, "bin"), filepath.Join(p.root, p.version.String(), "bin"))
		err = FileCopy(filepath.Join(tmp_build_dir, "bin"), filepath.Join(p.gvm.PkgsetRoot()))
		if err == nil {
			p.logger.Info(" * Installed binaries")
		}

		p.logger.Info("Installed", p.name, p.version.String())
	}
}
