package main

import "os"
import "path/filepath"

//import "io/ioutil"
//import "strings"

import . "source"
import . "version"

func (app *App) buildLocalPackage() (string, string) {
	wd, _ := os.Getwd()
	name := filepath.Base(wd)
	if app.pkgname != "" {
		name = app.pkgname
	}

	abspath, err := filepath.Abs(wd + "/..")
	if err != nil {
		app.Fatal("Failed to get parent folder")
	}

	return name, abspath
}

func (app *App) build() {
	name, path := app.buildLocalPackage()
	p := app.NewPackageDeprecated(name, nil, NewSource(path))
	app.Debug(p)
	p.Install(app.opts)
	return
}

func (app *App) clone() {
	name := app.readCommand()
	if name == "" {
		app.Fatal("Please specify package name")
	}
	found, version, source := app.FindPackageInCache(name, app.version)
	if found == false {
		app.Fatal("Couldn't find", name, app.version, "in any sources")
	}
	wd, _ := os.Getwd()
	source.Clone(name, version, wd)
	app.Message("Cloned", name, version)
	return
}

func (app *App) test() {
	name, path := app.buildLocalPackage()
	p := app.NewPackageDeprecated(name, nil, NewSource(path))
	app.Debug(p)
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

	p := app.NewPackageDeprecated(name, version, source)
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
	if name == "" {
		app.Fatal("Please specify package name")
	}
	found, version, source := app.FindPackageInCache(name, app.version)
	if found == false {
		app.Fatal(name, app.version, " not installed")
	}

	p := app.NewPackageDeprecated(name, version, source)
	app.Debug(p)
	p.Install(app.opts)
	app.StartDocServer(name)
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
