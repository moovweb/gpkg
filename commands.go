package main

import "os"
import "path/filepath"
//import "io/ioutil"
//import "strings"

func (app *App) build() {
	wd, _ := os.Getwd()
	name := filepath.Base(wd)
	if len(os.Args) > 1 {
		name = os.Args[1]
	}

	p := app.NewPackageFromSource(name, wd)
	app.Debug(p)
	//p.force_install = true
	p.Install("")
	return
}

func (app *App) install() {
	name := app.readCommand()
	if name == "" {
		app.Fatal("Please specify package name")
	}
	p := app.NewPackage(name, app.version)
	if p == nil {
		app.Fatal("Couldn't find", name, "in any sources")
	}
	app.Debug(p)
	p.Install("")
}

func (gpkg *App) sources() {
/*	if len(os.Args) < 2 {
		for _, src := range gpkg.gvm.sources {
			gpkg.logger.Info(src.root)
		}
		return
	}

	command := readCommand()
	if command == "add" {
		gpkg.gvm.AddSource(os.Args[1])
		gpkg.logger.Message("Added", os.Args[1], "to sources")
	} else if command == "remove" {
		if !gpkg.gvm.RemoveSource(os.Args[1]) {
			gpkg.logger.Error("Couldn't remove", os.Args[1])
		}
		gpkg.logger.Message("Removed", os.Args[1], "from sources")
	} else {
	}*/
}

func (gpkg *App) uninstall() {
/*	logger := gpkg.logger
	gvm := gpkg.gvm
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
	}*/
}

/*func (gpkg *App) listPackageVersions(pkg *Package) {
	versions := pkg.GetVersions()
	version_str := ""
	for n, version := range versions {
		version_str += version
		if n < len(versions) - 1 {
			version_str += ", "
		}
	}
	gpkg.logger.Info(pkg.name, "(" + version_str + ")")
}*/

/*func (gpkg *App) listPackage(pkg *Package) {
	logger := gpkg.logger
	logger.Message("Package Info:", pkg.name)
	logger.Info("  version:", pkg.version)
	data, err := ioutil.ReadFile(filepath.Join(pkg.root, pkg.version.String(), "manifest"))
	if err == nil {
		lines := strings.Split(string(data), "\n")
		for _, line := range lines {
			if len(line) > 7 && line[0:7] == ":source" {
				logger.Info("  source:", line[8:])
			}
		}
		deps := ""
		for _, line := range lines {
			if len(line) > 3 && line[0:3] == "pkg" {
				deps += "    " + line[4:] + "\n"
			}
		}
		if deps != "" {
			logger.Info("  deps:")
			logger.Info(deps)
		}
	}
}*/

func (app *App) list() {
/*	logger := gpkg.logger
	gvm := gpkg.gvm
	pkgname := readCommand()
	if pkgname != "" {
		var pkg *Package
		pkg = gvm.FindPackage(pkgname)
		if pkg != nil {
			version := flag.String("version", "", "Package version to install")
			flag.Parse()
			if *version != "" {
				pkg = gvm.FindPackageByVersion(pkgname, *version)
				if pkg != nil {
					gpkg.listPackage(pkg)
				} else {
					logger.Fatal("Package version not found")
				}
				return
			} else {
				logger.Message("\ngpkg list", pkg.name, "in", gvm.go_name + "@" + gvm.pkgset_name, "\n")
				gpkg.listPackageVersions(pkg)
			}
		} else {
			logger.Fatal("Package not found")
		}
	} else {
		logger.Message("\ngpkg package list", gvm.go_name + "@" + gvm.pkgset_name, "\n")
		pkgs := gvm.PackageList()
		for _, pkg := range pkgs {
			gpkg.listPackageVersions(pkg)
		}
		data, err := ioutil.ReadFile(filepath.Join(gvm.pkgset_root, "goinstall.log"))
		if err == nil {
			logger.Message("\ngoinstall package list", gvm.go_name + "@" + gvm.pkgset_name, "\n")
			logger.Info(strings.TrimSpace(string(data)))
		}
	}
	logger.Info()*/
}

