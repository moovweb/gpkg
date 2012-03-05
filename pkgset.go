package main

import "os"
import "path/filepath"
import "strings"
import "exec"

type Pkgset struct {
	gvm *Gvm
	logger *Logger
	pkgset *Pkgset
	g *Go

	root string
	name string
}

func (pkgset *Pkgset) FindPackageByVersion(name string, version string) *Package {
	_, err := os.Open(filepath.Join(pkgset.root, "pkg.gvm", name, version, "pkg"))
	if err == nil {
		return pkgset.NewPackage(name, version)
	}
	return nil
}

func (pkgset *Pkgset) InstallPackage(name string, version string) *Package {
	if pkgset.FindPackageByVersion(name, version) != nil {
		// TODO: Figure this out
		//gvm.logger.Fatal("Package", name, "already installed!")
		panic("Already installed")
	}
	p := pkgset.NewPackage(name, version)

	source := pkgset.gvm.FindSource(name)
	p.source = source	
	p.pkgset = pkgset
	p.Install()
	return p
}

func (pkgset *Pkgset) NewPackage(name string, version string) *Package {
	p := &Package{
		name: name,
		gvm: pkgset.gvm,
		logger: pkgset.logger,
		g: pkgset.g,
		pkgset: pkgset,
		version: version,
	}

	if version == "" {
		p.version = "0.0.src"
	}
	p.root = filepath.Join(pkgset.root, "pkg.gvm", p.name)
	return p
}

func (pkgset *Pkgset) PackageList() (pkglist[] *Package) {
	out, err := exec.Command("ls", filepath.Join(pkgset.root, "pkg.gvm")).CombinedOutput()
	if err == nil {
		pkgs := strings.Split(string(out), "\n")
		pkgs = pkgs[0:len(pkgs)-1]
		pkglist = make([]*Package, len(pkgs))
		for n, pkg := range pkgs {
			pkglist[n] = pkgset.NewPackage(pkg, "0.0.src")
		}
		return pkglist
	}
	return []*Package{}
}
