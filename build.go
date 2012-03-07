package main

import "os"
import "path/filepath"

func (gpkg *Gpkg) build() {
	gvm := gpkg.gvm
	wd, _ := os.Getwd()
	p := gvm.NewPackage(os.Args[1], "")
	p.source = wd
	p.Install()
	return
}
