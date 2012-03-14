package gpkg

import "io/ioutil"
import "strings"
import . "versions"

type Sources struct {
	list []Source
}

func NewSources(config_file string) (sources *Sources) {
	data, err := ioutil.ReadFile(config_file)
	if err != nil {
		return nil
	}

	sources = &Sources{}
	lines := strings.Split(string(data), "\n")
	count := 0
	for _, line := range lines {
		if line != "" && strings.TrimSpace(line)[0] != '#' {
			count++
		}
	}
/*	
	sources.list = make([]Source, count+1)
	sources.list[0] = NewGpkgSource("/home/jbussdieker/.gvm/pkgsets/release.r60.3/global/pkg.gvm")
	count = 1
*/
	sources.list = make([]Source, count)
	count = 0
	for _, line := range lines {
		if line != "" && strings.TrimSpace(line)[0] != '#' {
			source_location := strings.TrimSpace(line)
			sources.list[count] = NewSource(source_location)
			count++
		}
	}
	return
}

func (sources *Sources) Find(name string) *Version {
	return sources.FindBySpec(name, "*")
}

func (sources *Sources) FindBySpec(name string, spec string) *Version {
	for _, source := range sources.list {
		versions, src_err := source.Versions(name)
		if src_err != nil {
			continue
		}
		version, version_err := NewVersionFromMatch(versions, spec)
		if version_err == nil && version != nil {
			return version
		}
	}

	return nil
}

