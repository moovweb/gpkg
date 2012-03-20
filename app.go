package main

import "flag"
import . "gpkglib"
import . "pkg"

const VERSION = "0.1.19"

type App struct {
	*Gpkg
	args []string
	fs   *flag.FlagSet

	command string

	pkgname string
	version string
	opts    BuildOpts
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

func (app *App) addBuildFlags(local bool, build bool, test bool, install bool) {
	if local == false {
		app.fs.StringVar(&app.version, "version", "", "Package version to install")
	} else {
		app.fs.StringVar(&app.pkgname, "pkgname", "", "Name to give package being built. Default is the folder name.")
	}
	app.fs.BoolVar(&app.opts.Build, "build", build, "Build the package")
	app.fs.BoolVar(&app.opts.Test, "test", test, "Run package tests")
	app.fs.BoolVar(&app.opts.Install, "install", install, "Install the package")
}

func (app *App) readArgs() bool {
	app.command = app.readCommand()
	app.fs = flag.NewFlagSet("gpkg [command]", flag.ContinueOnError)
	log_level := app.fs.String("log", "info", "Log Level")
	if app.command == "install" || app.command == "uninstall" {
		app.addBuildFlags(false, true, true, true)
	}
	if app.command == "build" {
		app.addBuildFlags(true, true, true, true)
	}
	if app.command == "test" {
		app.addBuildFlags(true, true, true, false)
	}
	if app.command == "doc" {
		app.addBuildFlags(false, false, false, false)
	}
	err := app.fs.Parse(app.skipCommands())
	app.Gpkg = NewGpkg(*log_level)
	if err != nil {
		if app.command == "" {
			app.Info("Commands:")
			app.printUsage()
		}
		return false
	}
	return true
}

func (app *App) printUsage() {
	app.Info("  list      - List installed packages")
	app.Info("  doc       - Show documentation for a package")
	app.Info("  install   - Install a package")
	app.Info("  uninstall - Uninstall a package")
	app.Info("  empty     - Clear out all installed packages")
	app.Info("  build     - Build a package in the current directory")
	app.Info("  test      - Run tests on the package in the current directory")
	app.Info("  source    - List/Add/Remove sources for packages")
	app.Info("  version   - Print the gpkg version")
}

func (app *App) Close() {
	app.Gpkg.Close()
}

