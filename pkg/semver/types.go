package semver

const UNKNOWN_VERSION = "custombuild"

type Version struct {
	Major int64
	Minor int64
	Patch int64
}
