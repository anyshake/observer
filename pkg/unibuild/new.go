package unibuild

import (
	"strconv"
	"time"
)

func New(toolchainId, channel, commit, buildTime string) *UniBuild {
	if toolchainId == "" {
		toolchainId = "unspecified"
	}
	if buildTime == "" {
		buildTime = "0"
	}

	timestamp, _ := strconv.ParseInt(buildTime, 10, 64)
	return &UniBuild{
		toolchain:   matchToolchainById(toolchainId),
		toolchainId: toolchainId,
		channel:     channel,
		commit:      commit,
		time:        time.Unix(timestamp, 0).UTC(),
	}
}
