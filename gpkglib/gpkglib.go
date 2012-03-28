package gpkglib

import "path/filepath"
import "os"
import "io/ioutil"
import "strings"
import "strconv"
import "os/exec"

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

func (gpkg *Gpkg) NewPackageDeprecated(name string, version *Version, source Source) *PackageDeprecated {
	p := NewPackageDeprecated(gpkg.Gvm, name, version, filepath.Join(gpkg.PkgsetRoot(), "pkg.gvm", name), source, gpkg.tmpdir, gpkg.Logger)
	return p
}

func (gpkg *Gpkg) EmptyPackages() error {
	return os.RemoveAll(filepath.Join(gpkg.PkgsetRoot(), "pkg.gvm"))
}

func (gpkg *Gpkg) StartDocServer(name string) {
	gopath := filepath.Join(gpkg.tmpdir, name)
	os.Setenv("GOPATH", gopath)
	gpkg.Message("Starting documentation server...")
	gpkg.Debug("GOPATH is", gopath)
	gpkg.Info("http://localhost:6060/pkg/src")
	cmd := exec.Command("godoc", "-http", ":6060")
	cmd.Run()
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

func (gpkg *Gpkg) Close() {
	os.RemoveAll(gpkg.tmpdir)
}
