package main

import "os"
import "path/filepath"

func (gpkg *Gpkg) build() {
	gvm := gpkg.gvm
	wd, _ := os.Getwd()
	var p *Package
	if len(os.Args) > 1 {
		p = gvm.NewPackage(os.Args[1], "")
	} else {
		p = gvm.NewPackage(filepath.Base(wd), "")
	}
	p.source = wd
	p.Install(gpkg.tmpdir)
	os.RemoveAll(gpkg.tmpdir)
	return
}
