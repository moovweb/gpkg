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
import . "version"

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

func (gpkg *Gpkg) NewPackage(name string, version *Version, source Source) *Package {
	p := NewPackage(gpkg.Gvm, name, version, filepath.Join(gpkg.PkgsetRoot(), "pkg.gvm", name), source, gpkg.tmpdir, gpkg.Logger)
	return p
}

func (gpkg *Gpkg) EmptyPackages() os.Error {
	return os.RemoveAll(filepath.Join(gpkg.PkgsetRoot(), "pkg.gvm"))
}

func (gpkg *Gpkg) VersionList(name string) (list []*Version) {
	out, err := exec.Command("ls", filepath.Join(gpkg.PkgsetRoot(), "pkg.gvm", name)).CombinedOutput()
	if err == nil {
		versions := strings.Split(string(out), "\n")
		versions = versions[0 : len(versions)-1]
		list = make([]*Version, len(versions))
		for n, version_str := range versions {
			v := NewVersion(version_str)
			list[n] = v
		}
		return list
	}
	return []*Version{}
}

func (gpkg *Gpkg) GoinstallList() (list []string) {
	out, err := ioutil.ReadFile(filepath.Join(gpkg.PkgsetRoot(), "goinstall.log"))
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

func (gpkg *Gpkg) PackageList() (list []string) {
	out, err := exec.Command("ls", filepath.Join(gpkg.PkgsetRoot(), "pkg.gvm")).CombinedOutput()
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
