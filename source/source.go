package source

import "path/filepath"
import "os"

import . "version"

type SourceError struct{ msg string }

func NewSourceError(msg string) *SourceError { return &SourceError{msg: msg} }
func (e *SourceError) Error() string        { return "Source Error: " + e.msg }

type Source interface {
	String() string
	Clone(string, *Version, string) error
	Versions(string) ([]Version, error)
	List() []string
}

type PackageSource interface {
	Source
	Load(string, *Version, string) error
	Delete(string, *Version) error
}

func NewSource(root string) Source {
	if root[0] == '/' {
		return NewLocalSource(root)
	}
	return NewGitSource(root)
}

func cleanDest(dest string, name string) error {
	err := os.MkdirAll(dest, 0755)
	if err != nil {
		return NewSourceError(err.Error())
	}
	err = os.RemoveAll(filepath.Join(dest, name))
	if err != nil {
		return NewSourceError(err.Error())
	}
	return nil
}
