package main

import "path/filepath"

type Go struct {
	gvm *Gvm
	logger *Logger
	pkgset *Pkgset
	name string
	root string
}

func (g *Go) NewPkgset(name string) (pkgset *Pkgset) {
	pkgset = &Pkgset{}
	pkgset.name = name
	pkgset.gvm = g.gvm
	pkgset.logger = g.logger
	pkgset.g = g
	pkgset.root = filepath.Join(g.gvm.root, "pkgsets", g.name, name)
	return
}
