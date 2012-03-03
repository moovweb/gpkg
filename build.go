package main

import "os"
import "pkg"
import "path"

func (g *GpkgApp) build() {
	wd, _ := os.Getwd()
	p := pkg.NewPackage(wd, path.Base(wd))
	println()
	println("Commands")
	for _, cmd := range p.Commands {
		println("  ", cmd)
	}
	println()
	println("Packages")
	for _, pkg := range p.Packages {
		println("  ", pkg)
	}
	println()
}

