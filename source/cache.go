package source

import "os/exec"
import "strings"
import "os"
import "path/filepath"

import . "util"
import . "version"

type CacheSource struct {
	root string
}

func NewCacheSource(root string) PackageSource {
	s := CacheSource{root: root}
	return PackageSource(s)
}

func (s CacheSource) String() string {
	return s.root
}

func (s CacheSource) Delete(name string, version *Version) error {
	err := os.RemoveAll(filepath.Join(s.root, name, version.String()))
	if err == nil {
		list, err := s.Versions(name)
		if err != nil {
			return NewSourceError("Failed to check if for other versions\n" + err.Error())
		} else if len(list) == 0 {
			err := os.RemoveAll(filepath.Join(s.root, name))
			if err != nil {
				return NewSourceError("Failed to main folder after removing last package")
			}
		}
	}
	return nil
}

func (s CacheSource) List() (list []string) {
	out, err := exec.Command("ls", s.root).CombinedOutput()
	if err == nil {
		pkgs := strings.Split(string(out), "\n")
		pkgs = pkgs[0 : len(pkgs)-1]
		list = make([]string, len(pkgs))
		for n, pkg := range pkgs {
			list[n] = pkg
		}
		return list
	}
	return []string{}
}

func (s CacheSource) Load(name string, version *Version, dest string) error {
	cleanDest(dest, name)
	err := FileCopy(filepath.Join(s.root, name, version.String(), "pkg"), dest)
	if err != nil {
		return NewSourceError(err.Error())
	}

	return nil
}

func (s CacheSource) Clone(name string, version *Version, dest string) error {
	cleanDest(dest, name)
	err := FileCopy(filepath.Join(s.root, name, version.String(), "src", name), dest)
	if err != nil {
		return NewSourceError(err.Error())
	}

	return nil
}

func (s CacheSource) Versions(name string) (list []Version, err error) {
	out, oserr := exec.Command("ls", filepath.Join(s.root, name)).CombinedOutput()
	if err == nil {
		versions := strings.Split(string(out), "\n")
		versions = versions[0 : len(versions)-1]
		list = make([]Version, len(versions))
		for n, version_str := range versions {
			v := NewVersion(version_str)
			if v == nil {
				return []Version{}, NewSourceError("Failed to create version for install package!")
			}
			list[n] = *v
		}
		return list, nil
	}
	return []Version{}, NewSourceError(oserr.Error())
}
