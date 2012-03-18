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

func (app *App) uninstall() {
	pkgname := app.readCommand()
	if pkgname == "" {
		app.Fatal("Please specify package name")
	}
	if app.version == "" {
		p := app.FindPackage(pkgname)
		app.Info(p)
		if p != nil {
			if app.DeletePackages(pkgname) {
				app.Message("Deleted", pkgname)
			} else {
				app.Fatal("Couldn't delete", pkgname)
			}
		} else {
			app.Fatal("Invalid package name")
		}
	} else {
		p := app.FindPackageByVersion(pkgname, app.version)
		app.Info(p)
		if p != nil {
			if app.DeletePackage(pkgname, app.version) {
				app.Message("Deleted", pkgname, "version", app.version)
			} else {
				app.Fatal("Couldn't delete", pkgname, "version", app.version)
			}
		} else {
			app.Fatal("Invalid package name or version")
		}
	}
}

func (app *App) list() {
	app.Message("\ngpkg package list", app.Gvm.GoName + "@" + app.Gvm.PkgsetName, "\n")
	pkgs := app.PackageList()
	for _, pkg := range pkgs {
		output := pkg + " ("
		versions := app.VersionList(pkg)
		for n, version := range versions {
			output += version
			if n < len(versions) - 1 {
				output += ", "
			}
		}
		output += ")"
		app.Info(output)
	}
	app.Message("\ngoinstall package list", app.Gvm.GoName + "@" + app.Gvm.PkgsetName, "\n")
	pkgs = app.GoinstallList()
	for _, pkg := range pkgs {
		app.Info(pkg)
	}
	app.Info()
}

