package gpkg
/*
import "os"
import "path/filepath"

type FileSource struct {
	root string
}

func NewFileSource(root string) FileSource {
	return FileSource{root: root}
}

func (fs FileSource) Name() string {
	return fs.root
}

func (fs FileSource) Find(name string) *PackageNode {
	return fs.FindBySpec(name, "*")
}

func (fs FileSource) FindBySpec(name string, spec string) *PackageNode {
	src_dir := filepath.Join(fs.root, name)
	_, err := os.Open(src_dir)
	if err == nil {
		p := NewPackage(name, fs, NewVersion("0.0.0"))
		pp := NewPackageNode(p, "*", []Version{})
		return pp
	}
	return nil
}

func (fs FileSource) Clone(p *Package, dest string) Error {
	src_dir := filepath.Join(fs.root, p.name)
	dest_dir := filepath.Join(dest, p.name)
	_, err := os.Open(dest_dir)
	if err != nil {
		FileCopy(src_dir, dest_dir)
	} else {
		return err
	}
	return nil
}
*/
