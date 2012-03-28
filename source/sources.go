package source

import "io/ioutil"
import "strings"

import . "version"
import . "errors"

type Sources struct {
	cache   PackageSource
	sources map[string]Source
}

func NewSources(cache_path string, config_file string) *Sources {
	s := &Sources{}
	s.sources = map[string]Source{}
	s.cache = NewCacheSource(cache_path)
	if config_file != "" {
		data, err := ioutil.ReadFile(config_file)
		if err != nil {
			return s
		}
		src_list := strings.Split(string(data), "\n")
		for _, src := range src_list {
			if src != "" && strings.TrimSpace(src)[0] != '#' {
				s.sources[src] = NewSource(strings.TrimSpace(src))
			}
		}
	}
	return s
}

func (s *Sources) findSource(name string) (bool, []Version, Source) {
	for _, source := range s.sources {
		versions, err := source.Versions(name)
		if err == nil {
			return true, versions, source
		}
	}
	return false, nil, nil
}

func (s *Sources) Add(source string) Error {
	s.sources[source] = NewSource(source)
	return nil
}

func (s *Sources) LoadFromCache(name string, version *Version, dest string) Error {
	return s.cache.Load(name, version, dest)
}

func (s *Sources) FindInCache(name string, spec string) (found bool, version *Version, source Source) {
	versions, verr := s.cache.Versions(name)
	if verr != nil {
		return false, nil, nil
	}
	v, err := NewVersionFromMatch(versions, spec)
	if err != nil {
		return false, nil, nil
	}
	if v == nil {
		return false, nil, nil
	}
	return true, v, s.cache
}

func (s *Sources) FindInSources(name string, spec string) (found bool, version *Version, source Source) {
	found, versions, source := s.findSource(name)
	if found == false {
		return false, nil, nil
	}
	v, err := NewVersionFromMatch(versions, spec)
	if err != nil {
		return false, nil, nil
	}
	if v == nil {
		return false, nil, nil
	}
	return true, v, source
}
