package main

import "exec"
import "os"
import "flag"

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
		version := flag.String("version", "", "Package version to install")
		flag.Parse()
		if *version == "" {
			gvm.InstallPackage(pkgname)
		} else {
			gvm.InstallPackageByVersion(pkgname, *version)
		}
	} else if command == "uninstall" {
		pkgname := readCommand()
		if pkgname == "" {
			logger.Fatal("Please specify package name")
		}
		version := flag.String("version", "", "Package version to install")
		flag.Parse()
		if *version == "" {
			p := gvm.FindPackage(pkgname)
			if p != nil {
				if gvm.DeletePackages(p.name) {
					logger.Message("Deleted", p.name)
				} else {
					logger.Fatal("Couldn't delete", p.name)
				}
			} else {
				logger.Fatal("Invalid package name")
			}
		} else {
			p := gvm.FindPackageByVersion(pkgname, *version)
			if p != nil {
				if gvm.DeletePackage(p) {
					logger.Message("Deleted", p.name, "version", p.version)
				} else {
					logger.Fatal("Couldn't delete", p.name, "version", p.version)
				}
			} else {
				logger.Fatal("Invalid package version")
			}
		}
	} else if command == "list" {
		logger.Message("\ngpkg list", gvm.go_name + "@" + gvm.pkgset_name, "\n")
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

