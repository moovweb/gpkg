package main

import "flag"

func (gpkg *Gpkg) install() {
	logger := gpkg.logger
	gvm := gpkg.gvm
	pkgname := readCommand()
	if pkgname == "" {
		logger.Fatal("Please specify package name")
	}
	version := flag.String("version", "", "Package version to install")
	flag.Parse()
	if *version == "" {
		gvm.InstallPackage(pkgname)
	} else {
		gvm.InstallPackageByVersion(pkgname, *version)
	}
}
