package semver_test

import (
	"strings"
	"testing"

	"github.com/anyshake/observer/pkg/semver"
)

func TestVersionPrecedence(t *testing.T) {
	t.Parallel()
	ordered := []string{
		"0.9.9", "1.0.0-alpha", "1.0.0-alpha.1", "1.0.0-alpha.2",
		"1.0.0-alpha.10", "1.0.0-alpha.99999999999999999999",
		"1.0.0-alpha.100000000000000000000", "1.0.0-alpha.beta",
		"1.0.0-beta", "1.0.0-beta.2", "1.0.0-beta.11",
		"1.0.0-rc.1", "1.0.0-rc.2", "1.0.0-rc.10", "1.0.0",
		"1.0.1", "1.1.0", "2.0.0",
	}
	for i, left := range ordered {
		for j, right := range ordered {
			a, b := version(left), version(right)
			if a.LessThan(b) != (i < j) || a.GreaterThan(b) != (i > j) || a.Equal(b) != (i == j) ||
				a.LessThanOrEqual(b) != (i <= j) || a.GreaterThanOrEqual(b) != (i >= j) {
				t.Errorf("inconsistent version comparison: %s versus %s", left, right)
			}
		}
	}
	if !version("1.0.0-1").LessThan(version("1.0.0--")) {
		t.Error("numeric prerelease identifier must precede a nonnumeric identifier")
	}
}

func TestVersionMetadataAndCompatibility(t *testing.T) {
	t.Parallel()
	v := semver.New("1", "2", "3", "rc.1")
	if v.String() != "v1.2.3-rc.1" || v.GetMajor() != 1 || v.GetMinor() != 2 || v.GetPatch() != 3 || v.GetPreRelease() != "rc.1" || !v.IsPreRelease() {
		t.Fatalf("unexpected prerelease version: %v", v)
	}
	if got := semver.New("", "", "", "").String(); got != "<custom-version>" {
		t.Fatalf("custom version = %q", got)
	}
	tests := []struct {
		current string
		other   string
		want    bool
	}{
		{"1.2.3", "1.9.0", true},
		{"1.2.3", "2.0.0", false},
		{"1.2.3-rc.1", "1.2.3", false},
		{"0.0.0", "0.1.0", false},
	}
	for _, tt := range tests {
		if got := version(tt.current).IsCompatible(version(tt.other)); got != tt.want {
			t.Errorf("%s compatible with %s = %v, want %v", tt.current, tt.other, got, tt.want)
		}
	}
}

func version(s string) *semver.Version {
	core, pre, _ := strings.Cut(s, "-")
	parts := strings.Split(core, ".")
	return semver.New(parts[0], parts[1], parts[2], pre)
}
