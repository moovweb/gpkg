package main

import "os"
import "path/filepath"
//import "io/ioutil"
//import "strings"

func (app *App) build() {
	wd, _ := os.Getwd()
	name := filepath.Base(wd)
	if app.pkgname != "" {
		name = app.pkgname
	}

	abspath, err := filepath.Abs(wd + "/..")
	if err != nil {
		app.Fatal("Failed to get parent folder")
	}

	p := app.NewPackage(name, "", abspath)
	app.Debug(p)
	//p.force_install = true
	p.Install()
	return
}

func (app *App) install() {
	name := app.readCommand()
	if name == "" {
		app.Fatal("Please specify package name")
	}
	found, version, source := app.FindPackageInSources(name, app.version)
	if found == false {
		app.Fatal("Couldn't find", name, app.version, "in any sources")
	}

	p := app.NewPackage(name, version, source)
	app.Debug(p)
	p.Install()
}

func (app *App) source() {
	command := app.readCommand()
	if command == "add" {
		source := app.readCommand()
		app.AddSource(source)
		app.Message("Added", source, "to sources")
	} else if command == "remove" {
		source := app.readCommand()
		if !app.RemoveSource(source) {
			app.Error("Couldn't remove", source)
		}
		app.Message("Removed", source, "from sources")
	} else if command == "list" || command == "" {
		for _, src := range app.SourceList() {
			app.Info(src.Root())
		}
	} else {
		app.Fatal("Invalid source command (" + command + ").\nValid choices are: add, remove, and list")
	}
}

func (app *App) uninstall() {
	pkgname := app.readCommand()
	if pkgname == "" {
		app.Fatal("Please specify package name")
	}
	found, version, _ := app.FindPackageInCache(pkgname, "")
	if found == false {
		app.Fatal("Invalid package name")
	}
	found, version, _ = app.FindPackageInCache(pkgname, app.version)
	if found == false {
		app.Fatal("Invalid package version")
	}
	if app.version == "" {
		if app.DeletePackages(pkgname) {
			app.Message("Deleted", pkgname)
		} else {
			app.Fatal("Couldn't delete", pkgname)
		}
	} else {
		if app.DeletePackage(pkgname, version) {
			app.Message("Deleted", pkgname, "version", version)
		} else {
			app.Fatal("Couldn't delete", pkgname, "version", version)
		}
	}
}

func (app *App) list() {
	app.Message("\ngpkg package list", app.Gvm.GoName+"@"+app.Gvm.PkgsetName, "\n")
	pkgs := app.PackageList()
	for _, pkg := range pkgs {
		output := pkg + " ("
		versions := app.VersionList(pkg)
		for n, version := range versions {
			output += version
			if n < len(versions)-1 {
				output += ", "
			}
		}
		output += ")"
		app.Info(output)
	}
	app.Message("\ngoinstall package list", app.Gvm.GoName+"@"+app.Gvm.PkgsetName, "\n")
	pkgs = app.GoinstallList()
	for _, pkg := range pkgs {
		app.Info(pkg)
	}
	app.Info()
}
