package source

import "path/filepath"
import "exec"
import "os"
import "strings"
import "fmt"

import . "version"
import . "util"
import . "errors"

type SourceError struct{ msg string }

func NewSourceError(msg string) *SourceError { return &SourceError{msg: msg} }
func (e *SourceError) String() string        { return "Source Error: " + e.msg }

type Source interface {
	String() string
	Clone(string, *Version, string) Error
	Delete(string, *Version) Error
	Versions(string) ([]Version, Error)
	List() []string
}

func NewSource(root string) Source {
	if root[0] == '/' {
		return NewLocalSource(root)
	}
	return NewGitSource(root)
}
///////////////////////////////////////////
// Common
///////////////////////////////////////////
func cleanDest(dest string, name string) Error {
	err := os.MkdirAll(dest, 0755)
	if err != nil {
		return NewSourceError(err.String())
	}
	err = os.RemoveAll(filepath.Join(dest, name))
	if err != nil {
		return NewSourceError(err.String())
	}
	return nil
}

///////////////////////////////////////////
// GIT SOURCE
///////////////////////////////////////////
type GitSource struct {
	root string
}

const GIT_TAG_PREFIX = "refs/tags/"

func NewGitSource(root string) Source {
	s := GitSource{root: root}
	return Source(s)
}

func (s GitSource) String() string {
	return s.root
}

func (s GitSource) Delete(name string, version *Version) Error {
	panic("Not implemented!")
}

func (s GitSource) List() []string {
	panic("Not implemented!")
	return []string{}
}

func (s GitSource) Clone(name string, version *Version, dest string) Error {
	cleanDest(dest, name)
	src_repo := s.root + "/" + name + ".git"
	dest_dir := filepath.Join(dest, name)
	out, err := exec.Command("git", "clone", src_repo, dest_dir).CombinedOutput()
	if err != nil {
		return NewSourceError(err.String() + "\n" + string(out))
	}
	if version != nil {
		err := os.Chdir(dest_dir)
		if err != nil {
			return NewSourceError(fmt.Sprintln("Unable to chdir to checkout version", version, "of", name))
		}
		_, err = exec.Command("git", "checkout", version.String()).CombinedOutput()
		if err != nil {
			return NewSourceError(fmt.Sprintln("Invalid version:", version, "of", name, "specified"))
		}
	}
	return nil
}

func (s GitSource) Versions(name string) (list []Version, err Error) {
	out, oserr := exec.Command("git", "ls-remote", s.root+"/"+name).CombinedOutput()
	if oserr != nil {
		return nil, NewSourceError(oserr.String())
	}
	lines := strings.Split(string(out), "\n")
	versions := make([]Version, len(lines))
	index := 0
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) > 1 {
			if len(fields[1]) > len(GIT_TAG_PREFIX) && fields[1][:len(GIT_TAG_PREFIX)] == GIT_TAG_PREFIX {
				version := NewVersion(fields[1][len(GIT_TAG_PREFIX):])
				if version != nil {
					versions[index] = *version
					index++
				}
			}
		}
	}
	return versions[0:index], nil
}

///////////////////////////////////////////
// LOCAL SOURCE
///////////////////////////////////////////
type LocalSource struct {
	root string
}

func NewLocalSource(root string) Source {
	s := LocalSource{root: root}
	return Source(s)
}

func (s LocalSource) String() string {
	return s.root
}

func (s LocalSource) Delete(name string, version *Version) Error {
	panic("Not implemented!")
}

func (s LocalSource) List() []string {
	panic("Not implemented!")
}

func (s LocalSource) Clone(name string, version *Version, dest string) Error {
	cleanDest(dest, name)
	err := FileCopy(filepath.Join(s.root, name), filepath.Join(dest, name))
	// TODO: This is a hack to get jenkins working on multitarget installs folder name != project name
	//if s.name != filepath.Base(dest) {
	//p.logger.Debug(" * Rename", filepath.Join(tmp_src_dir, filepath.Base(p.source)), "to", filepath.Join(tmp_src_dir, p.name))
	//os.Rename(filepath.Join(tmp_src_dir, filepath.Base(p.source)), filepath.Join(tmp_src_dir, p.name))
	//return NewSourceError("TODO: Fix package rename at install")
	//}
	// END TODO
	if err != nil {
		return NewSourceError(err.String())
	}

	return nil
}

func (s LocalSource) Versions(name string) (list []Version, err Error) {
	// TODO: This assumes theres a test for NewVersion("0.0.0")!
	return []Version{*NewVersion("0.0.0")}, nil
}

///////////////////////////////////////////
// Cache SOURCE
///////////////////////////////////////////
type CacheSource struct {
	root string
}

func NewCacheSource(root string) Source {
	s := CacheSource{root: root}
	return Source(s)
}

func (s CacheSource) String() string {
	return s.root
}

func (s CacheSource) Delete(name string, version *Version) Error {
	err := os.RemoveAll(filepath.Join(s.root, name, version.String()))
	if err == nil {
		list, err := s.Versions(name)
		if err != nil {
			return NewSourceError("Failed to check if for other versions\n" + err.String())
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

func (s CacheSource) Clone(name string, version *Version, dest string) Error {
	cleanDest(dest, name)
	err := FileCopy(filepath.Join(s.root, name, version.String(), "src", name), dest)
	if err != nil {
		return NewSourceError(err.String())
	}

	return nil
}

func (s CacheSource) Versions(name string) (list []Version, err Error) {
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
	return []Version{}, NewSourceError(oserr.String())
}
