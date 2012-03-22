package container

import "path/filepath"
import "os"

import . "errors"

type SimpleContainer struct {
	root string
}

func NewSimpleContainer(root string) Container {
	return Container(SimpleContainer{root: root})
}

func (c SimpleContainer) String() string {
	return c.root
}

func (c SimpleContainer) SrcDir() string {
	return filepath.Join(c.root, "src")
}

func (c SimpleContainer) PkgDir() string {
	return filepath.Join(c.root, "pkg")
}

func (c SimpleContainer) BinDir() string {
	return filepath.Join(c.root, "bin")
}

func (c SimpleContainer) Empty() Error {
	err := os.RemoveAll(c.root)
	if err != nil {
		return err
	}

	// TODO: Do the right thing here!
	err = os.MkdirAll(c.root, 0775)
	if err != nil { return err }
	err = os.MkdirAll(c.SrcDir(), 0775)
	if err != nil { return err }
	err = os.MkdirAll(c.PkgDir(), 0775)
	if err != nil { return err }
	err = os.MkdirAll(c.BinDir(), 0775)
	return err
}
