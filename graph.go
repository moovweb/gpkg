package main

import "os"
import "pkg"
import "path"
import "path/filepath"

func findImport(name string, packages[] string) bool {
	for _, base_pkg := range packages {
		if base_pkg == name {
			return true
		}
	}
	return false
}

func (g *GpkgApp) graph() {
	gvm_path := os.Getenv("GVM_ROOT")
	gvm_go_name := os.Getenv("gvm_go_name")

	base_pkgs := pkg.NewPackage(filepath.Join(gvm_path, "gos", gvm_go_name, "src/pkg"), "")

	wd, _ := os.Getwd()
	p := pkg.NewPackage(wd, path.Base(wd))
	println("==", path.Base(wd), "==")
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
	println("Imports")
	for _, imp := range p.Imports {
		if !findImport(imp, base_pkgs.Packages) && !findImport(imp, p.Packages) {
			println("  ", imp)
		}
	}
	println()
}
