package gpkg

import . "versions"

type SourceError struct { msg string }
func NewSourceError(msg string) *SourceError { return &SourceError{msg:msg} }
func (e *SourceError) String() string { return "Source Error: " + e.msg }

type Source interface {
	Name() string
	Versions(name string) ([]Version, *SourceError)
	Clone(name string, dest string) *SourceError
	SetVersion(name string, dest string, version Version) *SourceError
}

func NewSource(location string) Source {
	if len(location) == 0 {
		return nil
	}

/*
	if location[0] == '/' {
		return NewFileSource(location)
	}
*/
	return NewGitSource(location)
}

