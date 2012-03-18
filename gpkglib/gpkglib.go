package gpkglib

import "path/filepath"
import "os"
import "io/ioutil"
import "exec"
import "strings"
import "strconv"

import "github.com/moovweb/versions"

import . "logger"
import . "gvm"
import . "pkg"

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
	gpkg.tmpdir = filepath.Join(os.Getenv("GOROOT"), "tmp", strconv.Itoa(os.Getpid()))
	return gpkg
}

func (gpkg *Gpkg) Close() {
	os.RemoveAll(gpkg.tmpdir)
}

func (gpkg *Gpkg) NewPackage(name string, tag string) *Package {
	found, source := gpkg.FindSource(name, tag)
	if found == false {
		return nil
	}
	p := NewPackage(gpkg.Gvm, name, tag, filepath.Join(gpkg.PkgsetRoot(), "pkg.gvm", name), source + "/" + name, gpkg.tmpdir, gpkg.Logger)
	return p
}

func (gpkg *Gpkg) NewPackageFromSource(name string, source string) *Package {
	p := NewPackage(gpkg.Gvm, name, "", filepath.Join(gpkg.PkgsetRoot(), "pkg.gvm", name), source, gpkg.tmpdir, gpkg.Logger)
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

func (gvm *Gpkg) FindPackage(name string) *Package {
	var p *Package
	var tag string

	gvm.Trace("name", name)
	_, err := os.Open(filepath.Join(gvm.PkgsetRoot(), "pkg.gvm", name))
	if err == nil {
		dirs, err := ioutil.ReadDir(filepath.Join(gvm.PkgsetRoot(), "pkg.gvm", name))
		if err != nil {
			panic("No versions")
		}
		for _, dir := range dirs {
			this_version, err := versions.NewVersion(dir.Name)
			if err != nil {
				gvm.Info("bad version1", dir.Name, err)
				continue
			}
			if p != nil {
				current_version, err := versions.NewVersion(tag)
				if err != nil {
					gvm.Info("bad version2", tag, err)
					continue
				}
				matched, err := this_version.Matches("> " + current_version.String())
				if err != nil {
					gvm.Info("bad match", tag, err)
					continue
				} else if matched == true {
					tag = dir.Name
					p = gvm.NewPackage(name, dir.Name)
				}
			} else {
				tag = dir.Name
				p = gvm.NewPackage(name, dir.Name)
			}
		}
	}
	return p
}

func (gvm *Gpkg) VersionList(name string) (list[] string) {
	out, err := exec.Command("ls", filepath.Join(gvm.PkgsetRoot(), "pkg.gvm", name)).CombinedOutput()
	if err == nil {
		pkgs := strings.Split(string(out), "\n")
		pkgs = pkgs[0:len(pkgs)-1]
		list = make([]string, len(pkgs))
		for n, pkg := range pkgs {
			list[n] = pkg
		}
		return list
	}
	return []string{}
}

func (gvm *Gpkg) GoinstallList() (list[] string) {
	out, err := ioutil.ReadFile(filepath.Join(gvm.PkgsetRoot(), "goinstall.log"))
	if err == nil {
		pkgs := strings.Split(string(out), "\n")
		pkgs = pkgs[0:len(pkgs)-1]
		list = make([]string, len(pkgs))
		for n, pkg := range pkgs {
			list[n] = pkg
		}
		return list
	}
	return []string{}
}

func (gvm *Gpkg) PackageList() (list[] string) {
	out, err := exec.Command("ls", filepath.Join(gvm.PkgsetRoot(), "pkg.gvm")).CombinedOutput()
	if err == nil {
		pkgs := strings.Split(string(out), "\n")
		pkgs = pkgs[0:len(pkgs)-1]
		list = make([]string, len(pkgs))
		for n, pkg := range pkgs {
			list[n] = pkg
		}
		return list
	}
	return []string{}
}

