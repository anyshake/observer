package semver

import "strings"

func (v *Version) Equal(ver *Version) bool {
	return v.major == ver.major && v.minor == ver.minor && v.patch == ver.patch && v.preRelease == ver.preRelease
}

func (v *Version) LessThan(ver *Version) bool {
	if v.major < ver.major {
		return true
	}
	if v.major > ver.major {
		return false
	}

	if v.minor < ver.minor {
		return true
	}
	if v.minor > ver.minor {
		return false
	}

	if v.patch < ver.patch {
		return true
	}
	if v.patch > ver.patch {
		return false
	}

	// Compare pre-release version if major, minor, patch are equal
	if v.preRelease != "" && ver.preRelease == "" {
		return true
	}
	if v.preRelease == "" && ver.preRelease != "" {
		return false
	}
	left, right := strings.Split(v.preRelease, "."), strings.Split(ver.preRelease, ".")
	for i := 0; i < len(left) && i < len(right); i++ {
		if left[i] == right[i] {
			continue
		}
		leftNumeric, rightNumeric := isNumericIdentifier(left[i]), isNumericIdentifier(right[i])
		if leftNumeric != rightNumeric {
			return leftNumeric
		}
		if leftNumeric && len(left[i]) != len(right[i]) {
			return len(left[i]) < len(right[i])
		}
		return left[i] < right[i]
	}
	return len(left) < len(right)
}

func isNumericIdentifier(s string) bool {
	if s == "" {
		return false
	}
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

func (v *Version) GreaterThan(ver *Version) bool {
	return !v.LessThan(ver) && !v.Equal(ver)
}

func (v *Version) LessThanOrEqual(ver *Version) bool {
	return v.LessThan(ver) || v.Equal(ver)
}

func (v *Version) GreaterThanOrEqual(ver *Version) bool {
	return v.GreaterThan(ver) || v.Equal(ver)
}

func (v *Version) IsPreRelease() bool {
	return v.preRelease != ""
}

func (v *Version) IsCompatible(ver *Version) bool {
	if v.preRelease != "" {
		return false
	}

	if v.major == 0 && v.minor == 0 && v.patch == 0 {
		return false
	}

	if v.major != ver.major {
		return false
	}

	return true
}
