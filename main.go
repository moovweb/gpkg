package main

import "exec"
import "os"
import "path/filepath"

type Gpkg struct {
	gvm *Gvm
	logger *Logger
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

func main() {
	logger := NewLogger("", INFO)
	gpkg := &Gpkg{logger: logger}
	gpkg.NewGvm()
	command := readCommand()

	if command == "install" {
		gpkg.install()
	} else if command == "uninstall" {
		gpkg.uninstall()
	} else if command == "list" {
		gpkg.list()
	} else if command == "sources" {
		gpkg.sources()
	} else if command == "graph" {
		gpkg.graph()
	} else {
		logger.Fatal("Invalid command. Please use: list, install or uninstall")
	}
}

