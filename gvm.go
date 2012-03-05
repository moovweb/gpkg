package main

import "os"
import "exec"
import "io/ioutil"
import "path/filepath"
import "strings"

type Gvm struct {
	root string
	g *Go
	sources []string
	logger *Logger
}

func NewGvm(logger *Logger) *Gvm {
	gvm := &Gvm{logger: logger}
	gvm.root = os.Getenv("GVM_ROOT")

	go_name := os.Getenv("gvm_go_name")
	if go_name != "" {
		gvm.g = gvm.NewGo(go_name)
	}
	
	data, err := ioutil.ReadFile(filepath.Join(gvm.root, "config", "sources"))
	if err != nil {
		panic(err)
	}

	gvm.sources = strings.Split(string(data), "\n")
	return gvm
}

func (gvm *Gvm) NewGo(name string) (g *Go) {
	g = &Go{}
	g.name = name
	g.gvm = gvm
	g.logger = gvm.logger
	g.root = filepath.Join(gvm.root, "gos", name)

	pkgset_name := os.Getenv("gvm_pkgset_name")
	if pkgset_name != "" {
		g.pkgset = g.NewPkgset(pkgset_name)
	}
	return
}

func (gvm *Gvm) FindSource(pkgname string) string {
	for _, source := range gvm.sources {
		src := source + "/" + pkgname
		if src[0] == '/' {
			_, err := os.Open(src)
			if err == nil {
				return src
			}
		} else {
			_, err := exec.Command("git", "ls-remote", src).CombinedOutput()
			if err == nil {
				return src
			}
		}
	}
	return ""
}
