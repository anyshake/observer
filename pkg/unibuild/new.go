package unibuild

import (
	"strconv"
	"time"
)

func New(toolchainId, channel, commit, buildTime string) *UniBuild {
	if toolchainId == "" {
		toolchainId = "unsupported"
	}
	if channel == "" {
		channel = "self-build"
	}
	if commit == "" {
		commit = "unknown"
	}
	if buildTime == "" {
		buildTime = "0"
	}

	timestamp, _ := strconv.ParseInt(buildTime, 10, 64)
	return &UniBuild{
		Toolchain:   getToolchainById(toolchainId),
		ToolchainId: toolchainId,
		Channel:     channel,
		Commit:      commit,
		Time:        time.Unix(timestamp, 0).UTC(),
	}
}
