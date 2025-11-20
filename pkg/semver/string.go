package semver

import "fmt"

func (v *Version) String() string {
	if v.Major == 0 && v.Minor == 0 && v.Patch == 0 {
		return UNKNOWN_VERSION
	}
	return fmt.Sprintf("v%d.%d.%d", v.Major, v.Minor, v.Patch)
}
