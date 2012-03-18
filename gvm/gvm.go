package gvm

import "os"
import "io/ioutil"
import "path/filepath"
import "strings"
import "exec"
import "github.com/moovweb/versions"

import . "logger"

type Gvm struct {
	root string
	go_name string
	go_root string
	pkgset_name string
	pkgset_root string
	sources []string
	logger *Logger
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

	return gvm
}

func (gvm *Gvm) PkgsetRoot() string {
	return gvm.pkgset_root
}

func (gvm *Gvm) Root() string {
	return gvm.root
}

func (gvm *Gvm) AddSource(src string) bool {
	for _, check_src := range gvm.sources {
		if check_src == src {
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
	gvm.sources = []string{}
	count := 0
	for _, src := range src_list {
		if src != "" && strings.TrimSpace(src)[0] != '#' {
			gvm.sources = append(make([]string, len(gvm.sources) + 1), gvm.sources...)
			gvm.sources[count] = strings.TrimSpace(src)
			count++
		}
	}
	return true
}

func (gvm *Gvm) FindPackageByVersion(name string, version string) (bool, string) {
	gvm.logger.Trace("name", name)
	gvm.logger.Trace("version", version)
	_, err := os.Open(filepath.Join(gvm.pkgset_root, "pkg.gvm", name, version))
	if err == nil {
		return true, filepath.Join(gvm.pkgset_root, "pkg.gvm", name, version)
	}
	return false, ""
}

func (gvm *Gvm) FindPackage(name string) (found bool, version string, source string) {
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
				gvm.logger.Info("bad version1", dir.Name, err)
				continue
			}
			if found == true {
				current_version, err := versions.NewVersion(version)
				if err != nil {
					gvm.logger.Info("bad version2", version, err)
					continue
				}
				matched, err := this_version.Matches("> " + current_version.String())
				if err != nil {
					gvm.logger.Info("bad match", version, err)
					continue
				} else if matched == true {
					version = dir.Name
					source = filepath.Join(gvm.pkgset_root, "pkg.gvm", name, dir.Name)
				}
			} else {
				found = true
				version = dir.Name
				source = filepath.Join(gvm.pkgset_root, "pkg.gvm", name, dir.Name)
			}
		}
	}
	return found, version, source
}

func (gvm *Gvm) FindSource(name string, version string) (bool, string) {
	for _, source := range gvm.sources {
		src := source + "/" + name
		if src[0] == '/' {
			_, err := os.Open(src)
			if err == nil {
				return true, source
			}
		} else {
			_, err := exec.Command("git", "ls-remote", src).CombinedOutput()
			if err == nil {
				return true, source
			}
		}
	}
	return false, ""
}
