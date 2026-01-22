package main

type ExitReason int

const (
	ExitInterrupt ExitReason = iota
	ExitRestart
	ExitError
)

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
