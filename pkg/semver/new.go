package semver

import "strconv"

func New(major, minor, patch, pre string) *Version {
	majorInt, _ := strconv.ParseInt(major, 10, 64)
	minorInt, _ := strconv.ParseInt(minor, 10, 64)
	patchInt, _ := strconv.ParseInt(patch, 10, 64)
	return &Version{
		Major:      majorInt,
		Minor:      minorInt,
		Patch:      patchInt,
		PreRelease: pre,
	}
}
