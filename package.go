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
	logger *Logger
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

func (p *Package) Get() bool {
	p.src = filepath.Join(p.root, p.version, "src")
	os.MkdirAll(p.src, 0775)
	if p.source[0] == '/' {
		p.logger.Debug(" * Copying", p.name)
		err := FileCopy(p.source, p.src)
		if err != nil {
			return false
		}
	} else {
		p.logger.Debug(" * Downloading", p.name)
		_, err := exec.Command("git", "clone", p.source, p.src + "/" + p.name).CombinedOutput()
		if err != nil {
			return false
		}
	}
	return true
}

func (p *Package) LoadImports() bool {
	data, err := ioutil.ReadFile(filepath.Join(p.src, p.name, "Package.gvm"))
	if err == nil {
		p.logger.Debug(" * Loading deps for", p.name)
		for _, line := range strings.Split(string(data), "\n") {
			if len(line) > 3 && line[0:3] == "pkg" {
				params := strings.Split(line, " ")
				dep := p.gvm.FindPackageByVersion(params[1], "0.0.src")
				if dep == nil {
					dep = p.gvm.InstallPackage(params[1], "0.0.src")
				}
				if dep == nil {
					p.logger.Fatal("ERROR: Couldn't find " + params[1] + " in any sources")
				}

				os.MkdirAll(p.tmpimp, 0775)
				err = FileCopy(filepath.Join(dep.root, p.version, "pkg"), p.tmpimp)
				if err != nil {
					p.logger.Fatal("ERROR: Couldn't load import: " + dep.name)
				}
			}
		}
	}
	return true
}

func (p *Package) Build() bool {
	p.tmpdir = fmt.Sprintf("%s/tmp/%s-%d/%s", p.gvm.root, p.name, os.Getpid(), "build")
	p.tmpimp = fmt.Sprintf("%s/tmp/%s-%d/%s", p.gvm.root, p.name, os.Getpid(), "import")

	if !p.LoadImports() {
		p.logger.Error("Failed to load imports")
		return false
	}

	p.logger.Debug(" * Building", p.name, p.version)

	os.Chdir(filepath.Join(p.src, p.name))
	os.Setenv("GOPATH", p.tmpdir + ":" + p.tmpimp)
	if os.Getenv("BUILD_NUMBER") == "" {
		os.Setenv("BUILD_NUMBER", "src")
	}

	out, err := exec.Command("make", "-f", "Makefile.gvm").CombinedOutput()
	if err != nil {
		p.logger.Error("Failed to build")
		p.logger.Error(string(out))
		return false
	}

	p.logger.Info("Installing", p.name, p.version + "...")

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
	p.logger.Debug("Starting install of", p.name)
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
