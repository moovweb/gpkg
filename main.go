// gpkg is the package manager for gvm
//
package main

import "os/exec"
import "os"

func FileCopy(src string, dst string) (err error) {
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

func main() {
	command := readCommand()
	logger := NewLogger("gpkg: ", INFO)
	gvm := NewGvm(logger)

	if command == "install" {
		pkgname := readCommand()
		if pkgname == "" {
			logger.Fatal("Please specify package name")
		}
		gvm.InstallPackage(pkgname, "0.0.src")
	} else if command == "list" {
		pkgs := gvm.PackageList()
		for _, pkg := range pkgs {
			logger.Info(pkg.name, "(" + pkg.version + ")")
		}
	} else if command == "graph" {
		graph()
	} else {
		logger.Fatal("Invalid command. Please use: list, install or uninstall")
	}
}
