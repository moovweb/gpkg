package source

import "exec"
import "strings"
import "fmt"
import "os"
import "path/filepath"

import . "errors"
import . "version"

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
	src_repo := s.root + "/" + name
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
