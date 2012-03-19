package main

import "os"
import "path/filepath"
//import "io/ioutil"
//import "strings"

import . "source"
import . "version"
import . "pkg"

func (app *App) buildLocalPackage() *Package {
	wd, _ := os.Getwd()
	name := filepath.Base(wd)
	if app.pkgname != "" {
		name = app.pkgname
	}

	abspath, err := filepath.Abs(wd + "/..")
	if err != nil {
		app.Fatal("Failed to get parent folder")
	}

	return app.NewPackage(name, nil, NewSource(abspath))
}

func (app *App) build() {
	p := app.buildLocalPackage()
	app.Debug(p)
	p.Install(app.opts)
	return
}

func (app *App) test() {
	p := app.buildLocalPackage()
	app.Debug(p)
	app.opts.Test = true
	app.opts.Install = false
	p.Install(app.opts)
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
	p.Install(app.opts)
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
			app.Info(src)
		}
	} else {
		app.Fatal("Invalid source command (" + command + ").\nValid choices are: add, remove, and list")
	}
}

func (app *App) packageWithLogging(name string) (string, *Version) {
	if name == "" {
		app.Fatal("Please specify package name")
	}
	found, version, _ := app.FindPackageInCache(name, "")
	if found == false {
		app.Fatal("Invalid package name")
	}
	found, version, _ = app.FindPackageInCache(name, app.version)
	if found == false {
		app.Fatal("Invalid package version")
	}
	return name, version
}

func (app *App) doc() {
	name := app.readCommand()
	pkgname, version := app.packageWithLogging(name)
	app.StartDocServer(pkgname, version)
}

func (app *App) uninstall() {
	name := app.readCommand()
	pkgname, version := app.packageWithLogging(name)
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
	app.Message("\ngpkg package list", app.Gvm.String(), "\n")
	pkgs := app.PackageList()
	for _, pkg := range pkgs {
		output := pkg + " ("
		versions := app.VersionList(pkg)
		for n, version := range versions {
			output += version.String()
			if n < len(versions)-1 {
				output += ", "
			}
		}
		output += ")"
		app.Info(output)
	}
	app.Message("\ngoinstall package list", app.Gvm.String(), "\n")
	pkgs = app.GoinstallList()
	for _, pkg := range pkgs {
		app.Info(pkg)
	}
	app.Info()
}
