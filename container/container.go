package container

type Container interface {
	String() string
	SrcDir() string
	PkgDir() string
	BinDir() string
	Empty() error
}
