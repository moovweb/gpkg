package version

import "strings"

import "github.com/moovweb/versions"

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

