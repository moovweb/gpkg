package pkg

import "os"
import "path/filepath"
import "strings"
import "go/parser"
import "go/token"

type Package struct {
	root string
	compiled bool
	Packages[] string
	Commands[] string
}

func filter(file *os.FileInfo) bool {
	if strings.HasSuffix(file.Name, ".go") && !file.IsDirectory() {
		return true
	}
	return false
}

func NewPackage(root string) (p *Package) {
	p = &Package{root: root}
	p.Packages = make([]string, 0, 2048)
	p.Commands = make([]string, 0, 256)
	p.parsePackage(p.root)
	return 
}

func (p *Package) addItem(items[] string, name string) ([] string) {
	items = items[0:len(items)+1]
	if (len(items) < cap(items)) {
		items[len(items)-1] = name
	} else {
		panic("Too many items")
	}
	return items
}

func (p *Package) parsePackage(path string) {
	cur_dir, err := os.Open(path)
	if err != nil {
		panic("Failed to open pkg dir")
	}
	cur_dirs, err := cur_dir.Readdir(0)
	if err != nil {
		panic("Failed to read dir")
	}

	fset := token.NewFileSet()
	pkg_list, err := parser.ParseDir(fset, path, filter, 0)
	if err != nil {
		panic("Couldn't parse files")
	}

	for _, pkg := range pkg_list {
		if pkg.Name == "main" {
			p.Commands = p.addItem(p.Commands, filepath.Base(path))
			continue
		}
		if len(path) > len(p.root) {
			p.Packages = p.addItem(p.Packages, path[len(p.root)+1:])
		} else {
			p.Packages = p.addItem(p.Packages, pkg.Name)
		}
	}

	for _, next_dir := range cur_dirs {
		if next_dir.IsDirectory() && next_dir.Name != "test" && next_dir.Name != "_obj" && next_dir.Name != "_bin" {
			p.parsePackage(filepath.Join(path, next_dir.Name))
		}
	}
}

