// gpkg is the package manager for gvm
//
package main

import "exec"
import "os"
import "path/filepath"
import "strconv"

type Gpkg struct {
	gvm *Gvm
	logger *Logger
	tmpdir string
}

func FileCopy(src string, dst string) (err os.Error) {
	_, err = exec.Command("cp", "-r", src, dst).CombinedOutput()
	return
}

func readCommand() string {
	if len(os.Args) < 2 {
		return ""
	}
	os.Args = os.Args[1:]
	return os.Args[0]
}

func (gpkg *Gpkg) NewGvm() *Gvm {
	gpkg.gvm = &Gvm{logger: gpkg.logger}
	gpkg.gvm.root = os.Getenv("GVM_ROOT")
	gpkg.gvm.go_name = os.Getenv("gvm_go_name")
	gpkg.gvm.go_root = filepath.Join(gpkg.gvm.root, "gos", gpkg.gvm.go_name)
	gpkg.gvm.pkgset_name = os.Getenv("gvm_pkgset_name")
	gpkg.gvm.pkgset_root = filepath.Join(gpkg.gvm.root, "pkgsets", gpkg.gvm.go_name, gpkg.gvm.pkgset_name)

	if !gpkg.gvm.ReadSources() {
		gpkg.gvm.logger.Fatal("Failed to read source list")
	}

	return gpkg.gvm
}

func (gpkg *Gpkg) printUsage() {
	logger := gpkg.logger
	logger.Info("  list      - List installed packages")
	logger.Info("  install   - Install a package")
	logger.Info("  uninstall - Uninstall a package")
	logger.Info("  empty     - Clear out all installed packages")
	logger.Info("  build     - Build a package in the current directory")
	logger.Info("  sources   - List/Add/Remove sources for packages")
	logger.Info("  graph     - Generate dot graph output using the current directory")
	logger.Info("  version   - Print the gpkg version")
}

func main() {
	logger := NewLogger("", DEBUG)
	gpkg := &Gpkg{logger: logger}
	gpkg.NewGvm()
	gpkg.tmpdir = filepath.Join(gpkg.gvm.root, "tmp", strconv.Itoa(os.Getpid()))
	defer func() {
		os.RemoveAll(gpkg.tmpdir)
	}()
	command := readCommand()

	if command == "install" {
		gpkg.install()
	} else if command == "debug" {
		logger.Info(gpkg.gvm.FindPackage("manhattan"))
		return
	} else if command == "uninstall" {
		gpkg.uninstall()
	} else if command == "empty" {
		os.RemoveAll(filepath.Join(gpkg.gvm.pkgset_root, "pkg.gvm"))
	} else if command == "build" {
		gpkg.build()
	} else if command == "list" {
		gpkg.list()
	} else if command == "sources" {
		gpkg.sources()
	} else if command == "graph" {
		gpkg.graph()
	} else if command == "version" {
		logger.Info("0.0.2")
	} else if command == "help" {
		logger.Message("The following commands are available:")
		gpkg.printUsage()
	} else {
		logger.Error("Invalid command. Please use one of the following:")
		gpkg.printUsage()
	}
}

