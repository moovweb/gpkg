package main

import "exec"
import "os"

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

func main() {
	command := readCommand()
	logger := NewLogger("", INFO)
	gvm := NewGvm(logger)

	if command == "install" {
		pkgname := readCommand()
		if pkgname == "" {
			logger.Fatal("Please specify package name")
		}
		gvm.InstallPackage(pkgname)
	} else if command == "list" {
		logger.Info("\ngvm packages", gvm.go_name + "@" + gvm.pkgset_name, "\n")
		pkgs := gvm.PackageList()
		for _, pkg := range pkgs {
			versions := pkg.GetVersions()
			version_str := ""
			for n, version := range versions {
				version_str += version
				if n < len(versions) - 1 {
					version_str += ", "
				}
			}
			logger.Info(pkg.name, "(" + version_str + ")")
		}
		logger.Info()
	} else if command == "graph" {
		graph()
	} else {
		logger.Fatal("Invalid command. Please use: list, install or uninstall")
	}
}

