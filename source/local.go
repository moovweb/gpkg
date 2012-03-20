package source

import "path/filepath"

import . "errors"
import . "version"
import . "util"

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
