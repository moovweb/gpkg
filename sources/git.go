package gpkg

import "path/filepath"
import "exec"
import "strings"
import "os"
import . "versions"

const TAG_PREFIX = "refs/tags/"

type GitSource struct {
	root string
}

func NewGitSource(root string) GitSource {
	return GitSource{root: root}
}

func (gs GitSource) Name() string {
	return gs.root
}

func (gs GitSource) Versions(name string) ([]Version, *SourceError) {
	src_repo := filepath.Join(gs.root, name)
	out, oserr := exec.Command("git", "ls-remote", src_repo).CombinedOutput()
	if oserr != nil {
		return nil, NewSourceError(oserr.String())
	}
	lines := strings.Split(string(out), "\n")
	versions := make([]Version, len(lines))
	index := 0
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) > 1 {
			if len(fields[1]) > len(TAG_PREFIX) && fields[1][:len(TAG_PREFIX)] == TAG_PREFIX {
				version := NewVersion(fields[1][len(TAG_PREFIX):])
				if version != nil {
					versions[index] = *version
					index++
				}
			}
		}
	}
	return versions[0:index], nil
}

func (gs GitSource) Clone(name string, dest string) (*SourceError) {
	src_repo := filepath.Join(gs.root, name)
	dest_dir := filepath.Join(dest, name)
	_, err := exec.Command("git", "clone", src_repo, dest_dir).CombinedOutput()
	if err != nil {
		return NewSourceError(err.String())
	}
	return nil
}

func (gs GitSource) SetVersion(name string, dest string, version Version) (*SourceError) {
	dest_dir := filepath.Join(dest, name)
	pushd, err := os.Getwd()
	if err != nil {
		return NewSourceError(err.String())
	}
	err = os.Chdir(dest_dir)
	if err != nil {
		return NewSourceError(err.String())
	}
	_, err = exec.Command("git", "checkout", version.String()).CombinedOutput()
	if err != nil {
		return NewSourceError(err.String())
	}
	err = os.Chdir(pushd)
	if err != nil {
		return NewSourceError(err.String())
	}
	return nil
}

