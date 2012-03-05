package main

import "flag"
import "os"
import "path/filepath"

func (gpkg *Gpkg) install() {
	logger := gpkg.logger
	gvm := gpkg.gvm
	pkgname := readCommand()
	if pkgname == "" {
		logger.Fatal("Please specify package name")
	} else if pkgname == "." {
		wd, _ := os.Getwd()
		p := gvm.NewPackage(filepath.Base(wd), "")
		p.source = wd
		p.Install()
		return
	}

	version := flag.String("version", "", "Package version to install")
	flag.Parse()
	if *version == "" {
		gvm.InstallPackage(pkgname)
	} else {
		gvm.InstallPackageByVersion(pkgname, *version)
	}
}
