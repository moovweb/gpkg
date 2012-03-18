// gpkg is the package manager for gvm
//
package main

import "os"

import . "gpkglib"
import . "logger"

const VERSION = "0.0.6"

type App struct {
	*Gpkg
}

func readCommand() string {
	if len(os.Args) < 2 {
		return ""
	}
	os.Args = os.Args[1:]
	return os.Args[0]
}

func (app *App) printUsage() {
	app.Info("Usage: gpkg [command]")
	app.Info()
	app.Info("Commands:")
	app.Info("  list      - List installed packages")
	app.Info("  install   - Install a package")
	app.Info("  uninstall - Uninstall a package")
	app.Info("  empty     - Clear out all installed packages")
	app.Info("  build     - Build a package in the current directory")
	app.Info("  sources   - List/Add/Remove sources for packages")
	app.Info("  version   - Print the gpkg version")
}

func main() {
	app := App{Gpkg:NewGpkg(DEBUG)}
	defer app.Gpkg.Close()
	command := readCommand()

	if command == "install" {
		app.install()
	} else if command == "debug" {
		//logger.Info(gpkg.gvm.FindPackage("manhattan"))
		return
	} else if command == "uninstall" {
		app.uninstall()
	} else if command == "empty" {
		//os.RemoveAll(filepath.Join(app.PkgsetRoot(), "pkg.gvm"))
	} else if command == "build" {
		app.build()
	} else if command == "list" {
		app.list()
	} else if command == "sources" {
		app.sources()
	} else if command == "version" {
		app.Info(VERSION)
	} else if command == "help" {
		app.Message("The following commands are available:")
		app.printUsage()
	} else {
		app.Error("Invalid command.")
		app.printUsage()
	}
}

