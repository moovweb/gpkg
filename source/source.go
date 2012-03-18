package source

import "path/filepath"
import "exec"
import "os"
import "strings"

import . "version"
import . "util"

type SourceError struct{ msg string }

func NewSourceError(msg string) *SourceError { return &SourceError{msg: msg} }
func (e *SourceError) String() string        { return "Source Error: " + e.msg }

type Source interface {
	Root() string
	Clone(string, string, string) *SourceError
	Versions(string) ([]Version, *SourceError)
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
func cleanDest(dest string, name string) *SourceError {
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

func (s GitSource) Root() string {
	return s.root
}

func (s GitSource) Clone(name string, version string, dest string) *SourceError {
	cleanDest(dest, name)
	src_repo := filepath.Join(s.root + "/" + name)
	dest_dir := filepath.Join(dest, name)
	_, err := exec.Command("git", "clone", src_repo, dest_dir).CombinedOutput()
	if err != nil {
		return NewSourceError(err.String())
	}
	return nil
}

func (s GitSource) Versions(name string) (list []Version, err *SourceError) {
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

func (s LocalSource) Root() string {
	return s.root
}

func (s LocalSource) Clone(name string, version string, dest string) *SourceError {
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

func (s LocalSource) Versions(name string) (list []Version, err *SourceError) {
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

func (s CacheSource) Root() string {
	return s.root
}

func (s CacheSource) Clone(name string, version string, dest string) *SourceError {
	cleanDest(dest, name)
	err := FileCopy(filepath.Join(s.root, name, version, "src", name), dest)
	if err != nil {
		return NewSourceError(err.String())
	}

	return nil
}

func (s CacheSource) Versions(name string) (list []Version, err *SourceError) {
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
