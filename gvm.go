package main

import "os"
import "io/ioutil"
import "path/filepath"
import "strings"
import "exec"
import "versions"

type Gvm struct {
	root string
	go_name string
	go_root string
	pkgset_name string
	pkgset_root string
	sources []*Source
	logger *Logger
}

func (gvm *Gvm) AddSource(src string) bool {
	for _, check_src := range gvm.sources {
		if check_src.root == src {
			gvm.logger.Fatal("Source already exists!")
		}
	}

	source_file := filepath.Join(gvm.root, "config", "sources")
	data, err := ioutil.ReadFile(source_file)
	if err != nil {
		return false
	}
	data = []byte(string(data) + "\n" + src)
	err = ioutil.WriteFile(source_file, data, 0644)
	if err != nil {
		return false
	}
	
	gvm.ReadSources()
	return true
}

func (gvm *Gvm) RemoveSource(src string) bool {
	source_file := filepath.Join(gvm.root, "config", "sources")
	data, err := ioutil.ReadFile(source_file)
	if err != nil {
		return false
	}
	src_list := strings.Split(string(data), "\n")
	output := ""
	for _, check_src := range src_list {
		if check_src != "" && strings.TrimSpace(check_src)[0] != '#' {
			if strings.TrimSpace(check_src) != src {
				output += check_src + "\n"
			}
		} else {
			output += check_src + "\n"
		}
	}
	err = ioutil.WriteFile(source_file, []byte(output), 0644)
	if err != nil {
		return false
	}
	return true
}

func (gvm *Gvm) ReadSources() bool {
	data, err := ioutil.ReadFile(filepath.Join(gvm.root, "config", "sources"))
	if err != nil {
		return false
	}
	src_list := strings.Split(string(data), "\n")
	count := 0
	for _, src := range src_list {
		if src != "" && strings.TrimSpace(src)[0] != '#' {
			count++
		}
	}
	gvm.sources = make([]*Source, count)
	count = 0
	for _, src := range src_list {
		if src != "" && strings.TrimSpace(src)[0] != '#' {
			gvm.sources[count] = NewSource(strings.TrimSpace(src))
			count++
		}
	}
	return true
}

func (gvm *Gvm) NewPackage(name string, tag string) *Package {
	p := &Package{
		gvm: gvm,
		logger: gvm.logger,
		name: name,
		tag: tag,
	}
	p.root = filepath.Join(p.gvm.pkgset_root, "pkg.gvm", p.name)
	return p
}

func (gvm *Gvm) InstallPackageByVersion(name string, tag string) *Package {
	p := gvm.NewPackage(name, tag)
	p.Install()
	return p
}

func (gvm *Gvm) InstallPackage(name string) *Package {
	gvm.logger.Trace("name", name)
	p := gvm.NewPackage(name, "")
	p.Install()
	return p
}

func (gvm *Gvm) FindPackageByVersion(name string, version string) *Package {
	gvm.logger.Trace("name", name)
	gvm.logger.Trace("version", version)
	_, err := os.Open(filepath.Join(gvm.pkgset_root, "pkg.gvm", name, version))
	if err == nil {
		p := gvm.NewPackage(name, version)
		return p
	}
	return nil
}

func (gvm *Gvm) DeletePackage(p *Package) bool {
	err := os.RemoveAll(filepath.Join(p.root, p.tag))
	if err == nil {
		if gvm.FindPackage(p.name) == nil {
			err := os.RemoveAll(filepath.Join(p.root))
			if err == nil {
				return true
			} else {
				return false
			}
		}
		return true
	}
	return false
}

func (gvm *Gvm) DeletePackages(name string) bool {
	err := os.RemoveAll(filepath.Join(gvm.pkgset_root, "pkg.gvm", name))
	if err == nil {
		return true
	}
	return false
}

func (gvm *Gvm) FindPackage(name string) *Package {
	var p *Package

	gvm.logger.Trace("name", name)
	_, err := os.Open(filepath.Join(gvm.pkgset_root, "pkg.gvm", name))
	if err == nil {
		dirs, err := ioutil.ReadDir(filepath.Join(gvm.pkgset_root, "pkg.gvm", name))
		if err != nil {
			panic("No versions")
		}
		for _, dir := range dirs {
			this_version, err := versions.NewVersion(dir.Name)
			if err != nil {
				gvm.logger.Info("bad version", dir.Name)
				continue
			}
			if p != nil {
				current_version, err := versions.NewVersion(p.tag)
				if err != nil {
					gvm.logger.Info("bad version", p.tag)
					continue
				}
				if this_version.NewerThan(current_version) {
					p = gvm.NewPackage(name, dir.Name)
				}
			} else {
				p = gvm.NewPackage(name, dir.Name)
			}
		}
	}
	return p
}

func (gvm *Gvm) PackageList() (pkglist[] *Package) {
	out, err := exec.Command("ls", filepath.Join(gvm.pkgset_root, "pkg.gvm")).CombinedOutput()
	if err == nil {
		pkgs := strings.Split(string(out), "\n")
		pkgs = pkgs[0:len(pkgs)-1]
		pkglist = make([]*Package, len(pkgs))
		for n, pkg := range pkgs {
			pkglist[n] = gvm.NewPackage(pkg, "")
		}
		return pkglist
	}
	return []*Package{}
}

