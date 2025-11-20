package semver

func (v *Version) Equal(ver *Version) bool {
	return v.Major == ver.Major && v.Minor == ver.Minor && v.Patch == ver.Patch && v.PreRelease == ver.PreRelease
}

func (v *Version) LessThan(ver *Version) bool {
	if v.Major < ver.Major {
		return true
	}
	if v.Major > ver.Major {
		return false
	}

	if v.Minor < ver.Minor {
		return true
	}
	if v.Minor > ver.Minor {
		return false
	}

	if v.Patch < ver.Patch {
		return true
	}
	if v.Patch > ver.Patch {
		return false
	}

	// Compare pre-release version if Major, Minor, Patch are equal
	if v.PreRelease != "" && ver.PreRelease == "" {
		return true
	}
	if v.PreRelease == "" && ver.PreRelease != "" {
		return false
	}
	return v.PreRelease < ver.PreRelease
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

func (v *Version) IsStable() bool {
	return v.PreRelease == ""
}
