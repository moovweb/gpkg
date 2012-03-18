package gpkglib

import "path/filepath"
import "os"
import "io/ioutil"
import "exec"
import "strings"
import "strconv"

//import "github.com/moovweb/versions"

import . "logger"
import . "gvm"
import . "pkg"
import . "source"

type Gpkg struct {
	*Gvm
	*Logger
	tmpdir string
}

func NewGpkg(loglevel string) *Gpkg {
	gpkg := &Gpkg{}
	gpkg.Logger = NewLogger("", LevelFromString(loglevel))
	gvm := NewGvm(gpkg.Logger)
	gpkg.Gvm = gvm
	gpkg.tmpdir = filepath.Join(os.Getenv("GVM_ROOT"), "tmp", strconv.Itoa(os.Getpid()))
	return gpkg
}

func (gpkg *Gpkg) NewPackage(name string, tag string) *Package {
	found, _, source := gpkg.FindSource(name, tag)
	if found == false {
		return nil
	}
	p := NewPackage(gpkg.Gvm, name, tag, filepath.Join(gpkg.PkgsetRoot(), "pkg.gvm", name), source, gpkg.tmpdir, gpkg.Logger)
	return p
}

func (gpkg *Gpkg) NewPackageFromSource(name string, source string) *Package {
	p := NewPackage(gpkg.Gvm, name, "", filepath.Join(gpkg.PkgsetRoot(), "pkg.gvm", name), NewSource(source), gpkg.tmpdir, gpkg.Logger)
	return p
}

func (gpkg *Gpkg) FindPackageByVersion(name string, version string) *Package {
	gpkg.Trace("name", name)
	gpkg.Trace("version", version)
	_, err := os.Open(filepath.Join(gpkg.PkgsetRoot(), "pkg.gvm", name, version))
	if err == nil {
		p := gpkg.NewPackage(name, version)
		return p
	}
	return nil
}

func (gvm *Gpkg) DeletePackage(name string, version string) bool {
	err := os.RemoveAll(filepath.Join(gvm.PkgsetRoot(), "pkg.gvm", name, version))
	if err == nil {
		if gvm.FindPackage(name) == nil {
			err := os.RemoveAll(filepath.Join(gvm.PkgsetRoot(), "pkg.gvm", name))
			if err == nil {
				return true
			} else {
				return false
			}
		}
		return true
	}
	return false
}

func (gvm *Gpkg) DeletePackages(name string) bool {
	err := os.RemoveAll(filepath.Join(gvm.PkgsetRoot(), "pkg.gvm", name))
	if err == nil {
		return true
	}
	return false
}

func (gvm *Gpkg) EmptyPackages() os.Error {
	return os.RemoveAll(filepath.Join(gvm.PkgsetRoot(), "pkg.gvm"))
}

func (gpkg *Gpkg) FindPackage(name string) *Package {
	found, version, source := gpkg.Gvm.FindPackage(name, "")
	if found != true {
		return nil
	}

	return NewPackage(gpkg.Gvm, name, version, filepath.Join(gpkg.PkgsetRoot(), "pkg.gvm", name), NewSource(source), filepath.Join(gpkg.PkgsetRoot(), "pkg.gvm", name), gpkg.Logger)
}

func (gvm *Gpkg) VersionList(name string) (list []string) {
	out, err := exec.Command("ls", filepath.Join(gvm.PkgsetRoot(), "pkg.gvm", name)).CombinedOutput()
	if err == nil {
		pkgs := strings.Split(string(out), "\n")
		pkgs = pkgs[0 : len(pkgs)-1]
		list = make([]string, len(pkgs))
		for n, pkg := range pkgs {
			list[n] = pkg
		}
		return list
	}
	return []string{}
}

func (gvm *Gpkg) GoinstallList() (list []string) {
	out, err := ioutil.ReadFile(filepath.Join(gvm.PkgsetRoot(), "goinstall.log"))
	if err == nil {
		pkgs := strings.Split(string(out), "\n")
		pkgs = pkgs[0 : len(pkgs)-1]
		list = make([]string, len(pkgs))
		for n, pkg := range pkgs {
			list[n] = pkg
		}
		return list
	}
	return []string{}
}

func (gvm *Gpkg) PackageList() (list []string) {
	out, err := exec.Command("ls", filepath.Join(gvm.PkgsetRoot(), "pkg.gvm")).CombinedOutput()
	if err == nil {
		pkgs := strings.Split(string(out), "\n")
		pkgs = pkgs[0 : len(pkgs)-1]
		list = make([]string, len(pkgs))
		for n, pkg := range pkgs {
			list[n] = pkg
		}
		return list
	}
	return []string{}
}

func (gpkg *Gpkg) Close() {
	os.RemoveAll(gpkg.tmpdir)
}
