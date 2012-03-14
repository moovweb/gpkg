package gpkg

import "github.com/moovweb/versions"

type VersionError struct { msg string }
func NewVersionError(msg string) *VersionError { return &VersionError{msg:msg} }
func (e *VersionError) String() string { return "Version Error: " + e.msg }

type Version struct {
	*versions.Version
}

func NewVersion(input string) (version *Version) {
	v, err := versions.NewVersion(input)
	if err == nil {
		version = &Version{}
		version.Version = v
	}
	return
}

func NewVersionFromMatch(versions[] Version, spec string) (*Version, *VersionError) {
	var version *Version
	for n, v := range versions {
		matched, err := v.Match(spec)
		if err != nil {
			return nil, err
		}
		if matched == true {
			if version != nil {
				if v.NewerThan(version) {
					version = &versions[n]
				}
			} else {
				version = &versions[n]
			}
		}
	}
	return version, nil
}

func (v *Version) Match(spec string) (bool, *VersionError) {
	if spec == "*" {
		return true, nil
	}
	match, err := v.Matches(spec)
	if err != nil {
		return false, NewVersionError(err.String())
	}

	if match == false {
		return false, nil
	}
	return true, nil
}

func (v *Version) NewerThan(target *Version) bool {
	matched, err := v.Matches("> " + target.String())
	if err != nil {
		return false
	}
	if matched == true {
		return true
	}
	return false
}
