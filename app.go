package main

import "flag"
import . "gpkglib"

const VERSION = "0.1.11"

type App struct {
	*Gpkg
	args    []string
	command string
	version string
	pkgname string
	fs      *flag.FlagSet
}

func NewApp(args []string) *App {
	app := &App{args: args}
	if !app.readArgs() {
		return nil
	}
	return app
}

func (app *App) readCommand() string {
	if len(app.args) < 2 {
		return ""
	}
	app.args = app.args[1:]
	return app.args[0]
}

func (app *App) skipCommands() []string {
	for n, arg := range app.args {
		if arg[0] == '-' {
			return app.args[n:]
		}
	}
	return app.args
}

func (app *App) readArgs() bool {
	app.command = app.readCommand()
	app.fs = flag.NewFlagSet("gpkg [command]", flag.ContinueOnError)
	log_level := app.fs.String("log", "DEBUG", "Log Level")
	if app.command == "install" || app.command == "uninstall" {
		app.fs.StringVar(&app.version, "version", "", "Package version to install")
	}
	if app.command == "build" {
		app.fs.StringVar(&app.pkgname, "pkgname", "", "Name to give package being built. Default is the folder name.")
	}
	err := app.fs.Parse(app.skipCommands())
	app.Gpkg = NewGpkg(*log_level)
	if err != nil {
		app.Info("Commands:")
		app.printUsage()
		return false
	}
	return true
}

func (app *App) printUsage() {
	app.Info("  list      - List installed packages")
	app.Info("  install   - Install a package")
	app.Info("  uninstall - Uninstall a package")
	app.Info("  empty     - Clear out all installed packages")
	app.Info("  build     - Build a package in the current directory")
	app.Info("  source    - List/Add/Remove sources for packages")
	app.Info("  version   - Print the gpkg version")
}
