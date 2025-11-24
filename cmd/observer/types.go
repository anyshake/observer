package main

const (
	officialBuildChannel = "github-actions-ci"
	startupDescription   = "Listen to the whispering earth."
)

var (
	versionMajor      string
	versionMinor      string
	versionPatch      string
	versionPreRelease string
)

var (
	buildTimestamp string
	buildToolchain string
	buildChannel   string
	buildCommit    string
)
