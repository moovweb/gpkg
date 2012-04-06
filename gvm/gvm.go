package gvm

import "os"
import "io/ioutil"
import "path/filepath"
import "strings"

import . "gpkg/version"

import . "gpkg/logger"
import . "gpkg/source"

type Gvm struct {
	go_name     string
	pkgset_name string

	root        string
	go_root     string
	pkgset_root string
	sources     []Source
	cache       PackageSource
	logger      *Logger
}

func NewGvm(logger *Logger) *Gvm {
	gvm := &Gvm{logger: logger}
	gvm.root = os.Getenv("GVM_ROOT")
	gvm.go_name = os.Getenv("gvm_go_name")
	gvm.go_root = filepath.Join(gvm.root, "gos", gvm.go_name)
	gvm.pkgset_name = os.Getenv("gvm_pkgset_name")
	gvm.pkgset_root = filepath.Join(gvm.root, "pkgsets", gvm.go_name, gvm.pkgset_name)

	if !gvm.ReadSources() {
		gvm.logger.Fatal("Failed to read source list")
	}
	gvm.cache = NewCacheSource(filepath.Join(gvm.root, "pkgsets", gvm.go_name, gvm.pkgset_name, "pkg.gvm"))

	return gvm
}

func (gvm *Gvm) String() string {
	return gvm.go_name + "@" + gvm.pkgset_name
}

func (gvm *Gvm) PkgsetRoot() string {
	return gvm.pkgset_root
}

func (gvm *Gvm) AddSource(src string) bool {
	for _, check_src := range gvm.sources {
		if check_src.String() == src {
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
	found := false
	for _, check_src := range src_list {
		if check_src != "" && strings.TrimSpace(check_src)[0] != '#' {
			if strings.TrimSpace(check_src) != src {
				output += check_src + "\n"
			} else {
				found = true
			}
		} else {
			output += check_src + "\n"
		}
	}
	if found == false {
		gvm.logger.Fatal("Source not found!")
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
	gvm.sources = make([]Source, count)
	count = 0
	for _, src := range src_list {
		if src != "" && strings.TrimSpace(src)[0] != '#' {
			gvm.sources[count] = NewSource(strings.TrimSpace(src))
			count++
		}
	}
	return true
}

func (gvm *Gvm) SourceList() []Source {
	return gvm.sources
}

func (gvm *Gvm) DeletePackage(name string, version *Version) bool {
	err := gvm.cache.Delete(name, version)
	if err != nil {
		return false
	}
	return true
}

func (gvm *Gvm) DeletePackages(name string) bool {
	versions, err := gvm.cache.Versions(name)
	if err != nil {
		return false
	}
	for _, version := range versions {
		err := gvm.cache.Delete(name, &version)
		if err != nil {
			return false
		}
	}
	return true
}

func (gvm *Gvm) PackageList() (list []string) {
	return gvm.cache.List()
}

func (gvm *Gvm) VersionList(name string) (list []Version) {
	versions, err := gvm.cache.Versions(name)
	if err != nil {
		return nil
	}
	return versions
}

func (gvm *Gvm) FindPackageInCache(name string, spec string) (found bool, version *Version, source Source) {
	versions, verr := gvm.cache.Versions(name)
	if verr != nil {
		return false, nil, nil
	}
	v, err := NewVersionFromMatch(versions, spec)
	if err != nil {
		return false, nil, nil
	}
	if v == nil {
		return false, nil, nil
	}
	return true, v, gvm.cache
}

func (gvm *Gvm) FindPackageInSources(name string, spec string) (found bool, version *Version, source Source) {
	found, versions, source := gvm.FindPackageSource(name)
	if found == false {
		return false, nil, nil
	}
	v, err := NewVersionFromMatch(versions, spec)
	if err != nil {
		return false, nil, nil
	}
	if v == nil {
		return false, nil, nil
	}
	return true, v, source
}

func (gvm *Gvm) FindPackageSource(name string) (bool, []Version, Source) {
	for _, source := range gvm.sources {
		versions, err := source.Versions(name)
		gvm.logger.Trace("FindSource: ", versions)
		if err == nil {
			return true, versions, source
		}
	}
	return false, nil, nil
}
