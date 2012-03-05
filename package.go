package main

import "exec"
import "os"
import "io/ioutil"
import "path/filepath"
import "strings"
import "fmt"

type Package struct {
	gvm *Gvm
	root string
	name string
	version string
	source string
	src string
	tmpdir string
	tmpimp string
	tmpsrc string
	logger *Logger
	deps map[string]*Package
}

func (p *Package) FindSource() bool {
	for _, source := range p.gvm.sources {
		src := source + "/" + p.name
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

func (p *Package) GetVersions() []string {
	dirs, _ := ioutil.ReadDir(p.root)
	versions := make([]string, len(dirs))
	for n, d := range dirs {
		versions[n] = d.Name
	}
	return versions
}

func (p *Package) Get() bool {
	p.tmpsrc = filepath.Join(p.gvm.root, "tmp", fmt.Sprintf("%d", os.Getpid()), p.name, "src")
	os.MkdirAll(p.tmpsrc, 0775)
	if p.source[0] == '/' {
		p.logger.Debug(" * Copying", p.name)
		err := FileCopy(p.source, p.tmpsrc)
		if err != nil {
			return false
		}
	} else {
		p.logger.Debug(" * Downloading", p.name)
		_, err := exec.Command("git", "clone", p.source, p.tmpsrc + "/" + p.name).CombinedOutput()
		if err != nil {
			return false
		}
	}

	if p.version != "" {
		p.logger.Debug(" * Checking out ", p.version)
		err := os.Chdir(p.tmpsrc + "/" + p.name)
		if err != nil {
			p.logger.Fatal("Unable to chdir to checkout version", p.version, "of", p.name)
		}
		_, err = exec.Command("git", "checkout", p.version).CombinedOutput()
		if err != nil {
			p.logger.Fatal("Invalid version:", p.version, "of", p.name, "specified")
		}
	}

	if p.version == "" {
		v, err := ioutil.ReadFile(filepath.Join(p.tmpsrc, p.name, "VERSION"))
		if err == nil {
			p.version = strings.TrimSpace(string(v))
		} else {
			p.version = "0.0"
		}
		if os.Getenv("BUILD_NUMBER") != "" {
			p.version += "." + os.Getenv("BUILD_NUMBER")
		} else {
			err := os.Chdir(p.tmpsrc + "/" + p.name)
			out, err := exec.Command("git", "describe", "--exact-match", "--tags", "--match", "*.*.*").CombinedOutput()
			if err == nil {
				p.version = strings.TrimSpace(string(out))
			} else {
				p.version += ".src"
			}
		}
	}	

	p.src = filepath.Join(p.root, p.version, "src")
	os.MkdirAll(filepath.Join(p.root, p.version), 0775)
	FileCopy(p.tmpsrc, p.src)

	return true
}

func (p *Package) LoadImport(dep *Package) {
	//p.logger.Trace("dep", fmt.Sprint(dep))
	if p.deps[dep.name] != nil {
		if p.deps[dep.name].version != dep.version {
			p.logger.Error("Version conflict!")
			p.logger.Fatal(dep.name, "version:", p.deps[dep.name].version, "Imported already for", p.name)
		}
	} else {
		for _, subdep := range dep.deps {
			p.LoadImport(subdep)
		}
		p.deps[dep.name] = dep
	}

	//p.deps += "pkg " + dep.name + " " + dep.version + "\n"
	os.MkdirAll(p.tmpimp, 0775)
	err := FileCopy(filepath.Join(dep.root, dep.version, "pkg"), p.tmpimp)
	if err != nil {
		p.logger.Fatal("ERROR: Couldn't load import: " + dep.name)
	}
}

func (p *Package) LoadImports() bool {
	data, err := ioutil.ReadFile(filepath.Join(p.src, p.name, "manifest"))
	if err != nil {
		data, err = ioutil.ReadFile(filepath.Join(p.src, p.name, "Package.gvm"))
		if err != nil {
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
					dep = p.gvm.InstallPackageByVersion(params[1], params[2])
				}
			} else {
				dep = p.gvm.FindPackage(params[1])
				if dep == nil {
					dep = p.gvm.InstallPackage(params[1])
				}
			}
			if dep == nil {
				p.logger.Fatal("ERROR: Couldn't find " + params[1] + " in any sources")
			}
			p.LoadImport(dep)
		}
	}
	return true
}

func (p *Package) Build() bool {
	p.tmpdir = fmt.Sprintf("%s/tmp/%d/%s/%s", p.gvm.root, os.Getpid(), p.name, "build")
	p.tmpimp = fmt.Sprintf("%s/tmp/%d/%s/%s", p.gvm.root, os.Getpid(), p.name, "import")

	if !p.LoadImports() {
		p.logger.Error("Failed to load imports")
		return false
	}

	p.logger.Debug(" * Building", p.name, p.version)

	os.Chdir(filepath.Join(p.src, p.name))
	os.Setenv("GOPATH", p.tmpdir + ":" + p.tmpimp)
	old_build_number := os.Getenv("BUILD_NUMBER")	
	os.Setenv("BUILD_NUMBER", p.version)
	out, err := exec.Command("make", "-f", "Makefile.gvm").CombinedOutput()
	if err != nil {
		p.logger.Error("Failed to build")
		p.logger.Error(string(out))
		return false
	}
	os.Setenv("BUILD_NUMBER", old_build_number)

	p.logger.Info("Installing", p.name + "-" + p.version + "...")

	//ioutil.WriteFile(filepath.Join(p.root, p.version, "manifest"), []byte(p.deps), 0664)
	err = FileCopy(filepath.Join(p.tmpdir, "pkg"), filepath.Join(p.root, p.version, "pkg"))
	if err != nil {
		return false
	}

	err = FileCopy(filepath.Join(p.tmpdir, "bin"), filepath.Join(p.gvm.pkgset_root))
	if err == nil {
		p.logger.Debug(" * Installed binaries")
	}

	return true
}

func (p *Package) Install() {
	p.logger.Debug("Starting install of", p.name, p.version)
	if !p.FindSource() {
		p.logger.Fatal("ERROR Finding package")
	}
	if !p.Get() {
		p.logger.Fatal("ERROR Getting package source")
	}
	if !p.Build() {
		p.logger.Fatal("ERROR Building package")
	}
}
