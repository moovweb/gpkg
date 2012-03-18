package main

import "flag"
import . "gpkglib"

const VERSION = "0.1.7"

type App struct {
	*Gpkg
	args[] string
	command string
	version string
	fs *flag.FlagSet
}

func NewApp(args[] string) *App {
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
	app.fs = flag.NewFlagSet("gpkg", flag.ContinueOnError)
	log_level := app.fs.String("log", "DEBUG", "Log Level")
	if app.command == "install" {
		app.fs.StringVar(&app.version, "version", "", "Package version to install")
	}
	err := app.fs.Parse(app.skipCommands())
	if err != nil {
		return false
	}
	app.Gpkg = NewGpkg(*log_level)
	return true
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

