package source

import "path/filepath"
import "os"

import . "version"
import . "errors"

type SourceError struct{ msg string }

func NewSourceError(msg string) *SourceError { return &SourceError{msg: msg} }
func (e *SourceError) String() string        { return "Source Error: " + e.msg }

type Source interface {
	String() string
	Clone(string, *Version, string) Error
	Versions(string) ([]Version, Error)
	List() []string
}

type PackageSource interface {
	Source
	Delete(string, *Version) Error
}

func NewSource(root string) Source {
	if root[0] == '/' {
		return NewLocalSource(root)
	}
	return NewGitSource(root)
}

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

