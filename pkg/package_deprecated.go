package pkg

import "os"
import "io/ioutil"
import "path/filepath"
import "strings"
import "fmt"

import . "gvm"
import . "logger"
import . "source"
import . "version"
import . "specs"
import . "tools"
import . "util"
import . "container"

type BuildOptsDeprecated struct {
	Build       bool
	Test        bool
	Install     bool
	InstallDeps bool
	UseSystem   bool
}

type PackageDeprecated struct {
	Source
	BuildOptsDeprecated
	gvm     *Gvm
	root    string
	name    string
	version *Version
	tmpdir  string
	logger  *Logger
	deps    map[string]*PackageDeprecated
	specs   *Specs
	tool    Tool
}

func NewPackageDeprecated(gvm *Gvm, name string, version *Version, root string, Source Source, tmpdir string, logger *Logger) *PackageDeprecated {
	p := &PackageDeprecated{
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

func (p *PackageDeprecated) String() string {
	return fmt.Sprintln("   root:", p.root) +
		fmt.Sprintln("   name:", p.name) +
		fmt.Sprintln("version:", p.version) +
		fmt.Sprintln(" source:", p.Source) +
		fmt.Sprintln(" tmpdir:", p.tmpdir)
}

func (p *PackageDeprecated) Clone() os.Error {
	p.logger.Info("Downloading", p.name, p.version)
	source_container := NewSimpleContainer(filepath.Join(p.tmpdir, p.name))
	err := p.Source.Clone(p.name, p.version, source_container.SrcDir())
	if err != nil {
		return err
	}

	if p.version == nil {
		version_str := ""
		v, err := ioutil.ReadFile(filepath.Join(source_container.SrcDir(), p.name, "VERSION"))
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

func (p *PackageDeprecated) LoadImport(dep *PackageDeprecated, dir string) {
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

func (p *PackageDeprecated) LoadImports(dir string) bool {
	tmp_src_dir := filepath.Join(p.tmpdir, p.name, "src")
	specs, err := NewSpecs(filepath.Join(tmp_src_dir, p.name))
	if err != nil {
		p.logger.Fatal("Failed to load specs", err)
	}
	if specs == nil {
		p.logger.Debug(" * No dependencies found")
		p.specs = NewBlankSpecs(p.Source)
		return true
	} else {
		p.specs = specs
	}

	p.deps = map[string]*PackageDeprecated{}

	p.logger.Info(" * Loading dependencies for", p.name)
	for name, spec := range p.specs.List {
		var dep *PackageDeprecated
		found, version, source := p.gvm.FindPackageInCache(name, spec)
		if found == true {
			dep = NewPackageDeprecated(p.gvm, name, version, filepath.Join(p.gvm.PkgsetRoot(), "pkg.gvm", name, version.String()), source, p.tmpdir, p.logger)
		} else if p.BuildOptsDeprecated.Install == true {
			found, version, source := p.gvm.FindPackageInSources(name, spec)
			if found == true {
				dep = NewPackageDeprecated(p.gvm, name, version, filepath.Join(p.gvm.PkgsetRoot(), "pkg.gvm", name), source, p.tmpdir, p.logger)
				p.logger.Trace(dep)
				dep.Install(p.BuildOptsDeprecated)
				dep.root = filepath.Join(p.gvm.PkgsetRoot(), "pkg.gvm", name, version.String())
			}
		} else {
			p.logger.Fatal("ERROR: Package", name, spec, "not found locally and install=false")
		}

		if dep == nil {
			p.logger.Fatal("ERROR: Couldn't find " + name + " " + spec + " in any sources")
		}
		p.logger.Info("    -", dep.name, dep.version, "(Spec:", spec+")")
		p.LoadImport(dep, dir)
	}
	return true
}

func (p *PackageDeprecated) WriteManifest() {
	p.specs.Origin = p.Source
	manifest := p.specs.String()
	err := ioutil.WriteFile(filepath.Join(p.root, p.version.String(), "manifest"), []byte(manifest), 0664)
	if err != nil {
		p.logger.Fatal("Failed to write manifest file")
	}
}

func (p *PackageDeprecated) PrettyLog(buf string) (formatted string) {
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

func (p *PackageDeprecated) Build() bool {
	tmp_build_dir := filepath.Join(p.tmpdir, p.name, "build")
	tmp_import_dir := filepath.Join(p.tmpdir, p.name, "import")

	if !p.LoadImports(tmp_import_dir) {
		p.logger.Error("Failed to load imports")
		return false
	}

	if p.BuildOptsDeprecated.Build == true {
		p.logger.Info(" * Building")

		old_gopath := os.Getenv("GOPATH")
		gopath := tmp_build_dir + ":" + tmp_import_dir
		if p.BuildOptsDeprecated.UseSystem && old_gopath != "" {
			gopath = gopath + ":" + old_gopath
		} else {
		}
		os.Setenv("GOPATH", gopath)

		old_build_number := os.Getenv("BUILD_NUMBER")
		os.Setenv("BUILD_NUMBER", p.version.String())

		// Clean
		out, berr := p.tool.Clean()
		if berr != nil {
			p.logger.Trace(berr)
			p.logger.Trace(out)
			//return false
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
		if p.BuildOptsDeprecated.Test == true {
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

func (p *PackageDeprecated) Install(b BuildOptsDeprecated) {
	p.BuildOptsDeprecated = b

	build_container := NewSimpleContainer(filepath.Join(p.tmpdir, p.name, "build"))
	source_container := NewSimpleContainer(filepath.Join(p.tmpdir, p.name))

	err := p.Clone()
	if err != nil {
		p.logger.Fatal("ERROR Getting package source", err)
	}
	p.tool, err = NewTool(filepath.Join(source_container.SrcDir(), p.name))
	if !p.Build() {
		p.logger.Fatal("ERROR Building package")
	}

	// INSTALL
	//////////////////////
	if p.BuildOptsDeprecated.Install == true {
		install_container := NewSimpleContainer(filepath.Join(p.root, p.version.String()))
		p.logger.Info(" * Installing", p.name+"-"+p.version.String()+"...")
		err := install_container.Empty()
		if err != nil {
			p.logger.Fatal("Failed to remove old version")
		}
		p.WriteManifest()
		err = FileCopy(source_container.SrcDir(), install_container.String())
		if err != nil {
			p.logger.Fatal("Failed to copy source to install folder", err)
		}

		err = FileCopy(build_container.PkgDir(), install_container.String())
		if err != nil {
			p.logger.Fatal("Failed to copy libraries to install folder\n", err.String())
		}

		err = FileCopy(build_container.BinDir(), install_container.BinDir())
		err = FileCopy(build_container.BinDir(), filepath.Join(p.gvm.PkgsetRoot()))
		if err == nil {
			p.logger.Info(" * Installed binaries")
		}

		p.logger.Info("Installed", p.name, p.version.String())
	}
}
