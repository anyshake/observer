package uniarch

import (
	"strings"

	"github.com/blang/semver/v4"
)

func GetGoSdkVersion(version string) (semver.Version, error) {
	version = strings.TrimPrefix(version, "go")

	switch len(strings.Split(version, ".")) {
	case 1:
		version += ".0.0"
	case 2:
		version += ".0"
	}

	return semver.Parse(version)
}
