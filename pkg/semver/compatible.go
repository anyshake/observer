package semver

func (v *Version) CompatibleWith(ver *Version) bool {
	if v.String() == UNKNOWN_VERSION || ver.String() == UNKNOWN_VERSION {
		return false
	}
	if v.Major != ver.Major {
		return false
	}
	if v.Minor != ver.Minor {
		return false
	}
	return true
}
