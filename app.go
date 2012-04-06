package main

import "flag"
import "runtime"
import . "gpkg/gpkglib"
import . "gpkg/pkg"

const VERSION = "0.1.23"

type App struct {
	*Gpkg
	args []string
	fs   *flag.FlagSet

	command string

	pkgname string
	version string
	opts    BuildOptsDeprecated
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

func (app *App) addBuildFlags(local bool, build bool, test bool, install bool, system bool) {
	if local == false {
		app.fs.StringVar(&app.version, "version", "", "Package version to install")
	} else {
		app.fs.StringVar(&app.pkgname, "pkgname", "", "Name to give package being built. Default is the folder name.")
	}
	if build == true {
		app.fs.BoolVar(&app.opts.Build, "build", build, "Build the package")
		app.fs.BoolVar(&app.opts.Test, "test", test, "Run package tests")
		app.fs.BoolVar(&app.opts.Install, "install", install, "Install the package")
		app.fs.BoolVar(&app.opts.UseSystem, "system", system, "Include normal GOPATH during build")
		app.fs.StringVar(&app.opts.TargetOS, "target_os", runtime.GOOS, "Specify the target OS")
		app.fs.StringVar(&app.opts.TargetArch, "target_arch", runtime.GOARCH, "Specify the target architechure")
	}
}

func (app *App) readArgs() bool {
	app.command = app.readCommand()
	app.fs = flag.NewFlagSet("gpkg [command]", flag.ContinueOnError)
	log_level := app.fs.String("log", "info", "Log Level")
	if app.command == "install" || app.command == "uninstall" {
		app.addBuildFlags(false, true, true, true, true)
	}
	if app.command == "build" {
		app.addBuildFlags(true, true, true, true, true)
	}
	if app.command == "test" {
		app.addBuildFlags(true, true, true, false, true)
	}
	if app.command == "clone" {
		app.addBuildFlags(false, false, false, false, false)
	}
	if app.command == "doc" {
		app.addBuildFlags(false, false, false, false, false)
	}
	if app.command == "bundle" {
		app.addBuildFlags(false, false, false, false, false)
	}
	if app.command == "goget" {
		app.addBuildFlags(false, true, false, false, false)
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
	app.Info("  clone     - Clone the source from an installed package")
	app.Info("  build     - Build a package in the current directory")
	app.Info("  goget     - Run go get in a gpkg context")
	app.Info("  bundle    - Print manifest of dependent versions")
	app.Info("  test      - Run tests on the package in the current directory")
	app.Info("  source    - List/Add/Remove sources for packages")
	app.Info("  version   - Print the gpkg version")
}

func (app *App) Close() {
	app.Gpkg.Close()
}
