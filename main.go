package main

import "os"
import "flag"

type GpkgApp struct {
	command string
}

func showCommands() {
	println("Please choose a command:")
	println("  Choices are: build, graph, install, uninstall")
	os.Exit(1)
}

func New() (g *GpkgApp) {
	g = &GpkgApp{}
	g.readArgs()
	return
}

func (g *GpkgApp) readArgs() {
	if len(os.Args) > 1 {
		g.command = os.Args[1]
		if g.command != "build" && g.command != "graph" && g.command != "install" && g.command != "uninstall" {
			showCommands()
		}
		os.Args = os.Args[1:]
	} else {
		showCommands()
	}

	flag.Parse()
}

func (g *GpkgApp) start() {
	if g.command == "build" {
		g.build()
	} else if g.command == "install" {
		g.install()
	} else if g.command == "uninstall" {
		g.uninstall()
	} else if g.command == "list" {
		g.list()
	} else if g.command == "graph" {
		g.graph()
	}
}

func main() {
	g := New()
	g.start()
}
