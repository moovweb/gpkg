package gpkg
/*
import "path/filepath"
import "io/ioutil"

type GpkgSource struct {
	root string
}

func NewGpkgSource(root string) GpkgSource {
	return GpkgSource{root:root}
}

func (gps GpkgSource) Name() string {
	return gps.root
}

func (gps GpkgSource) readVersionList(name string) (data string) {
	dirs, err := ioutil.ReadDir(filepath.Join(gps.root, name))
	if err == nil {
		for _, dir := range dirs {
			data += dir.Name + "\n"
		}
	}
	return
}

func (gps GpkgSource) Find(name string) *PackageNode {
	return gps.FindBySpec(name, "*")
}

func (gps GpkgSource) FindBySpec(name string, spec string) *PackageNode {
	data := gps.readVersionList(name)
	version, versions := findMatches(data, 0, "", spec)
	if version != nil {
		p := NewPackage(name, gps, version)
		pn := NewPackageNode(p, spec, versions)
		return pn
	}
	return nil
}

func (gps GpkgSource) Clone(p *Package, dest string) Error {
	dest_dir := filepath.Join(dest, p.name)

	src := filepath.Join(gps.root, p.name, p.version.String(), "src", p.name)
	FileCopy(src, dest_dir)
	return nil
}
*/
