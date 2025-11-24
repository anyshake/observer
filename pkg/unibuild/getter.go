package unibuild

import "time"

func (build *UniBuild) GetToolchainId() string {
	return build.toolchainId
}

func (build *UniBuild) GetToolchain() *Toolchain {
	return build.toolchain
}

func (build *UniBuild) GetCommit() string {
	if build.commit == "" {
		return "<out-of-tree>"
	}
	return build.commit
}

func (build *UniBuild) GetChannel() string {
	if build.channel == "" {
		return "self-build"
	}
	return build.channel
}

func (build *UniBuild) GetTime() time.Time {
	return build.time
}
