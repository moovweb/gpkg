package main

import "os"
import "pkg"
import "path"
import "path/filepath"
import "strings"
import "fmt"

func findImport(name string, packages[] string) bool {
	for _, base_pkg := range packages {
		if base_pkg == name {
			return true
		}
	}
	return false
}

func graph() {
	gvm_path := os.Getenv("GVM_ROOT")
	gvm_go_name := os.Getenv("gvm_go_name")

	base_pkgs := pkg.NewPackage(filepath.Join(gvm_path, "gos", gvm_go_name, "src/pkg"), "")

	wd, _ := os.Getwd()
	p := pkg.NewPackage(wd, path.Base(wd))
	dep := make(map[string]string, 256)
	for _, imp := range p.Imports {
		if !findImport(imp, base_pkgs.Packages) && !findImport(imp, p.Packages) && path.Base(wd) != imp {
			dep[strings.Split(imp, "/")[0]] = strings.Split(imp, "/")[0]
		}
	}

	for _, str := range dep {
		fmt.Print(path.Base(wd), "->", str + ";")
	}
}
