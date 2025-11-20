package unibuild

import "time"

type Toolchain struct {
	Name   string
	GOOS   string
	GOARM  string
	GOMIPS string
	GOARCH string
}

type UniBuild struct {
	Toolchain   *Toolchain
	ToolchainId string
	Commit      string
	Channel     string
	Time        time.Time
}
