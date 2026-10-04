package unibuild_test

import (
	"testing"
	"time"

	"github.com/anyshake/observer/pkg/unibuild"
)

func TestBuildMetadata(t *testing.T) {
	t.Parallel()
	build := unibuild.New("linux_arm32_v7a", "stable", "abc123", "1700000000")
	if build.GetToolchainId() != "linux_arm32_v7a" || build.GetChannel() != "stable" || build.GetCommit() != "abc123" || !build.GetTime().Equal(time.Unix(1700000000, 0).UTC()) {
		t.Fatalf("unexpected build metadata: %#v", build)
	}
	want := unibuild.Toolchain{Name: "linux-arm32-v7a", GOOS: "linux", GOARCH: "arm", GOARM: "7"}
	if toolchain := build.GetToolchain(); toolchain == nil || *toolchain != want {
		t.Fatalf("toolchain = %#v, want %#v", toolchain, want)
	}
}

func TestBuildMetadataFallbacks(t *testing.T) {
	t.Parallel()
	build := unibuild.New("", "", "", "invalid")
	if build.GetToolchainId() != "unspecified" || build.GetToolchain() != nil || build.GetChannel() != "self-build" || build.GetCommit() != "<out-of-tree>" || !build.GetTime().Equal(time.Unix(0, 0).UTC()) {
		t.Fatalf("fallback build metadata = %#v", build)
	}
	unknown := unibuild.New("unknown-toolchain", "dev", "commit", "")
	if unknown.GetToolchain() != nil || !unknown.GetTime().Equal(time.Unix(0, 0).UTC()) {
		t.Fatalf("unknown build metadata = %#v", unknown)
	}
}
