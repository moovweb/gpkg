package gpkg

import "path/filepath"
import "strconv"
import "os"

type Container struct {
	root string
}

func NewContainer() *Container {
	root := filepath.Join(os.TempDir(), "gpkg-container" + strconv.Itoa(os.Getpid()))
	return NewContainerWithRoot(root)
}

func NewContainerWithRoot(root string) *Container {
	os.MkdirAll(filepath.Join(root, "src"), 0755)
	return &Container{root:root}
}

func (container *Container) Close() {
	os.RemoveAll(container.root)
}

func (container *Container) SrcDir() string {
	return filepath.Join(container.root, "src")
}

func (container *Container) Root() string {
	return container.root
}

