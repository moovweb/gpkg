package container

import . "errors"

type Container interface {
	String() string
	SrcDir() string
	PkgDir() string
	BinDir() string
	Empty() Error
}

