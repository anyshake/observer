package semver

import "fmt"

func (v *Version) String() string {
	if v.Major == 0 && v.Minor == 0 && v.Patch == 0 {
		return UNKNOWN_VERSION
	}
	base := fmt.Sprintf("v%d.%d.%d", v.Major, v.Minor, v.Patch)
	if v.PreRelease != "" {
		return fmt.Sprintf("%s-%s", base, v.PreRelease)
	}
	return base
}
