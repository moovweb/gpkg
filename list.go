package main

func (gpkg *Gpkg) listPackageVersions(pkg *Package) {
	versions := pkg.GetVersions()
	version_str := ""
	for n, version := range versions {
		version_str += version
		if n < len(versions) - 1 {
			version_str += ", "
		}
	}
	gpkg.logger.Info(pkg.name, "(" + version_str + ")")
}

func (gpkg *Gpkg) list() {
	logger := gpkg.logger
	gvm := gpkg.gvm
	pkgname := readCommand()
	if pkgname != "" {
		pkg := gvm.FindPackage(pkgname)
		gpkg.listPackageVersions(pkg)
	} else {
		logger.Message("\ngpkg list", gvm.go_name + "@" + gvm.pkgset_name, "\n")
		pkgs := gvm.PackageList()
		for _, pkg := range pkgs {
			gpkg.listPackageVersions(pkg)
		}
		logger.Info()
	}
}
