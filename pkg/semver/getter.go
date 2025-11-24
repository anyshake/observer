package semver

import "fmt"

func (v *Version) String() string {
	if v.major == 0 && v.minor == 0 && v.patch == 0 {
		return "<custom-version>"
	}
	base := fmt.Sprintf("v%d.%d.%d", v.major, v.minor, v.patch)
	if v.preRelease != "" {
		return fmt.Sprintf("%s-%s", base, v.preRelease)
	}
	return base
}

func (v *Version) GetMajor() int64 {
	return v.major
}

func (v *Version) GetMinor() int64 {
	return v.minor
}

func (v *Version) GetPatch() int64 {
	return v.patch
}

func (v *Version) GetPreRelease() string {
	return v.preRelease
}
