package main

import "os"
import "github.com/moovweb/gpkg/pkg"

func main() {
	wd, _ := os.Getwd()
	p := pkg.NewPackage(wd)
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
