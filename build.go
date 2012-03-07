package main

import "os"
import "path/filepath"

func (gpkg *Gpkg) build() {
	gvm := gpkg.gvm

	wd, _ := os.Getwd()
	p := gvm.NewPackage(filepath.Base(wd), "")
	p.source = wd
	p.Install()
	return
}
