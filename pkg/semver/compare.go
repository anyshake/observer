package semver

func (v *Version) Equal(ver *Version) bool {
	return v.String() == ver.String()
}

func (v *Version) LessThan(ver *Version) bool {
	if v.Major > ver.Major {
		return false
	}
	if v.Minor > ver.Minor {
		return false
	}
	return v.Patch < ver.Patch
}

func (v *Version) GreaterThan(ver *Version) bool {
	if v.Major < ver.Major {
		return false
	}
	if v.Minor < ver.Minor {
		return false
	}
	return v.Patch > ver.Patch
}

func (v *Version) LessThanOrEqual(ver *Version) bool {
	return v.LessThan(ver) || v.Equal(ver)
}

func (v *Version) GreaterThanOrEqual(ver *Version) bool {
	return v.GreaterThan(ver) || v.Equal(ver)
}
