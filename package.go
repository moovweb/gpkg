package main

import "exec"
import "os"
import "io/ioutil"
import "path/filepath"
import "strings"
import "github.com/moovweb/gpkg/versions"

type Package struct {
	gvm *Gvm
	root string
	name string
	tag string
	version *versions.Version
	source string
	src string
	tmpdir string
	logger *Logger
	deps map[string]*Package
}

func (p *Package) String() string {
	return "   root: " + p.root + "\n" +
		"   name: " + p.name + "\n" +
		"    tag: " + p.tag + "\n" +
		" source: " + p.source + "\n" +
		"    src: " + p.src + "\n"
}

func (p *Package) GetVersions() []string {
	dirs, _ := ioutil.ReadDir(p.root)
	versions := make([]string, len(dirs))
	for n, d := range dirs {
		versions[n] = d.Name
	}
	return versions
}

func (p *Package) FindSource() bool {
	for _, source := range p.gvm.sources {
		src := source.root + "/" + p.name
		if src[0] == '/' {
			_, err := os.Open(src)
			if err == nil {
				p.source = src
				return true
			}
		} else {
			_, err := exec.Command("git", "ls-remote", src).CombinedOutput()
			if err == nil {
				p.source = src
				return true
			}
		}
	}
	return false
}

func (p *Package) Get() bool {
	tmp_src_dir := filepath.Join(p.tmpdir, p.name, "src")
	os.MkdirAll(tmp_src_dir, 0775)
	if p.source[0] == '/' {
		p.logger.Debug(" * Copying", p.name)
		err := FileCopy(p.source, tmp_src_dir)
		// TODO: This is a hack to get jenkins working on multitarget installs folder name != project name
		if p.name != filepath.Base(p.source) {
			p.logger.Debug("Rename", filepath.Join(tmp_src_dir, filepath.Base(p.source)), "to", filepath.Join(tmp_src_dir, p.name))
			os.Rename(filepath.Join(tmp_src_dir, filepath.Base(p.source)), filepath.Join(tmp_src_dir, p.name))
		}
		// END TODO
		if err != nil {
			return false
		}
	} else {
		p.logger.Debug(" * Downloading", p.name)
		_, err := exec.Command("git", "clone", p.source, tmp_src_dir + "/" + p.name).CombinedOutput()
		if err != nil {
			return false
		}
	}

	if p.tag != "" {
		p.logger.Debug(" * Checking out ", p.tag)
		err := os.Chdir(tmp_src_dir + "/" + p.name)
		if err != nil {
			p.logger.Fatal("Unable to chdir to checkout version", p.tag, "of", p.name)
		}
		_, err = exec.Command("git", "checkout", p.tag).CombinedOutput()
		if err != nil {
			p.logger.Fatal("Invalid version:", p.tag, "of", p.name, "specified")
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
			err := os.Chdir(tmp_src_dir + "/" + p.name)
			out, err := exec.Command("git", "describe", "--exact-match", "--tags", "--match", "*.*.*").CombinedOutput()
			if err == nil {
				p.tag = strings.TrimSpace(string(out))
			} else {
				if p.tag != "" {
					p.tag += ".src"
				} else {
					p.tag = "src"
				}
			}
		}
	}	

	p.src = filepath.Join(p.root, p.tag, "src")
	os.MkdirAll(filepath.Join(p.root, p.tag), 0775)
	FileCopy(tmp_src_dir, p.src)

	return true
}

func (p *Package) LoadImport(dep *Package, dir string) {
	//p.logger.Trace("dep", fmt.Sprint(dep))
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
	err := FileCopy(filepath.Join(dep.root, dep.tag, "pkg"), dir)
	if err != nil {
		p.logger.Fatal("ERROR: Couldn't load import: " + dep.name)
	}
}

func (p *Package) LoadImports(dir string) bool {
	data, err := ioutil.ReadFile(filepath.Join(p.src, p.name, "manifest"))
	if err != nil {
		data, err = ioutil.ReadFile(filepath.Join(p.src, p.name, "Package.gvm"))
		if err != nil {
			p.logger.Debug("No dependencies found")
			return true
		}
	}

	p.deps = make(map[string]*Package, 64)

	p.logger.Debug(" * Loading deps for", p.name)
	for _, line := range strings.Split(string(data), "\n") {
		if len(line) > 3 && line[0:3] == "pkg" {
			params := strings.Split(line, " ")
			var dep *Package
			if len(params) > 2 {
				dep = p.gvm.FindPackageByVersion(params[1], params[2])
				if dep == nil {
					dep = &Package{name:params[1],tag:params[2]}
					dep.Install(p.tmpdir)
				}
			} else {
				dep = p.gvm.FindPackage(params[1])
				if dep == nil {
					dep = &Package{name:params[1]}
					dep.Install(p.tmpdir)
				}
			}
			if dep == nil {
				p.logger.Fatal("ERROR: Couldn't find " + params[1] + " in any sources")
			}
			p.LoadImport(dep, dir)
		}
	}
	return true
}

func (p *Package) WriteManifest() {
	manifest := ":source " + p.source + "\n"
	for _, pkg := range p.deps {
		manifest += "pkg " + pkg.name + " " + pkg.tag + "\n"
	}
	ioutil.WriteFile(filepath.Join(p.root, p.tag, "manifest"), []byte(manifest), 0664)
}

func (p *Package) Build() bool {
	tmp_build_dir := filepath.Join(p.tmpdir, p.name, "build")
	tmp_import_dir := filepath.Join(p.tmpdir, p.name, "import")

	if !p.LoadImports(tmp_import_dir) {
		p.logger.Error("Failed to load imports")
		return false
	}

	p.logger.Debug(" * Building", p.name, p.tag)

	os.Chdir(filepath.Join(p.src, p.name))
	os.Setenv("GOPATH", tmp_build_dir + ":" + tmp_import_dir)
	old_build_number := os.Getenv("BUILD_NUMBER")	
	os.Setenv("BUILD_NUMBER", p.tag)
	_, err := os.Open("Makefile.gvm")
	if err == nil {
		out, err := exec.Command("make", "-f", "Makefile.gvm").CombinedOutput()
		if err != nil {
			p.logger.Error("Failed to build with Makefile.gvm")
			p.logger.Error(string(out))
			return false
		} else {
			p.logger.Debug(string(out))
		}
	} else {
		out, err := exec.Command("gb", "-bi").CombinedOutput()
		if err != nil {
			p.logger.Error("Failed to build with gb")
			p.logger.Error(string(out))
			return false
		} else {
			p.logger.Debug(string(out))
		}
	}

	p.WriteManifest()
	
	os.Setenv("BUILD_NUMBER", old_build_number)

	p.logger.Info("Installing", p.name + "-" + p.tag + "...")

	err = FileCopy(filepath.Join(tmp_build_dir, "pkg"), filepath.Join(p.root, p.tag, "pkg"))
	if err != nil {
		return false
	}

	err = FileCopy(filepath.Join(tmp_build_dir, "bin"), filepath.Join(p.gvm.pkgset_root))
	if err == nil {
		p.logger.Debug(" * Installed binaries")
	}

	return true
}

func (p *Package) Install(tmpdir string) {
	p.tmpdir = tmpdir
	p.logger.Debug("Starting install of", p.name, p.tag)
	if p.source == "" {
		if !p.FindSource() {
			p.logger.Fatal("ERROR Finding package")
		}
	}
	if !p.Get() {
		p.logger.Fatal("ERROR Getting package source")
	}
	if !p.Build() {
		p.gvm.DeletePackage(p)
		p.logger.Fatal("ERROR Building package")
	}
}
