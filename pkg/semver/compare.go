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

func (v *Version) MoreThan(ver *Version) bool {
	if v.Major < ver.Major {
		return false
	}
	if v.Minor < ver.Minor {
		return false
	}
	return v.Patch > ver.Patch
}
