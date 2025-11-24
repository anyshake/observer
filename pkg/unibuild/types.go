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
	toolchain   *Toolchain
	toolchainId string
	commit      string
	channel     string
	time        time.Time
}
