package main

import "flag"

func (gpkg *Gpkg) uninstall() {
	logger := gpkg.logger
	gvm := gpkg.gvm
	pkgname := readCommand()
	if pkgname == "" {
		logger.Fatal("Please specify package name")
	}
	version := flag.String("version", "", "Package version to install")
	flag.Parse()
	if *version == "" {
		p := gvm.FindPackage(pkgname)
		if p != nil {
			if gvm.DeletePackages(p.name) {
				logger.Message("Deleted", p.name)
			} else {
				logger.Fatal("Couldn't delete", p.name)
			}
		} else {
			logger.Fatal("Invalid package name")
		}
	} else {
		p := gvm.FindPackageByVersion(pkgname, *version)
		if p != nil {
			if gvm.DeletePackage(p) {
				logger.Message("Deleted", p.name, "version", p.version)
			} else {
				logger.Fatal("Couldn't delete", p.name, "version", p.version)
			}
		} else {
			logger.Fatal("Invalid package version")
		}
	}
}

