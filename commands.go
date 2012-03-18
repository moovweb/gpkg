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

	p := app.NewPackageFromSource(name, abspath)
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
	if app.version == "" {
		p := app.FindPackage(pkgname)
		if p != nil {
			app.Trace(p)
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
		if p != nil {
			app.Trace(p)
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
