package version

import "strings"

import "github.com/moovweb/versions"

type VersionError struct{ msg string }

func NewVersionError(msg string) *VersionError { return &VersionError{msg: msg} }
func (e *VersionError) String() string         { return "Version Error: " + e.msg }

type Version struct {
	versions.Version
}

func NewVersion(version string) *Version {
	v, err := versions.NewVersion(version)
	if err != nil {
		return nil
	}
	return &Version{Version: *v}
}

func (v *Version) String() string {
	version_str := v.Version.String()
	parts := strings.Split(version_str, ".")
	if len(parts) == 3 && parts[2] == "0" {
		version_str = parts[0] + "." + parts[1] + ".src"
	}
	return version_str
}

func NewVersionFromMatch(versions []Version, spec string) (*Version, *VersionError) {
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
		return false, NewVersionError(err.Error())
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
