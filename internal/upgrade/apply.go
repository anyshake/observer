//go:build !windows
// +build !windows

package upgrade

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/anyshake/observer/pkg/semver"
)

func (u *Helper) ApplyUpgrade(version *semver.Version, release []byte) error {
	if version == nil {
		return fmt.Errorf("version is nil")
	}
	if u.appliedVer != nil {
		if version.Equal(u.appliedVer) {
			return nil
		}
	}

	// Check for write permission
	dir := filepath.Dir(u.currentExePath)
	probe, err := os.CreateTemp(dir, ".upgrade-perm-*")
	if err != nil {
		return err
	}
	probe.Close()
	_ = os.Remove(probe.Name())

	ts := time.Now().UTC().Format("20060102150405")
	tmp := fmt.Sprintf("%s.%s.new", u.currentExePath, ts)
	old := fmt.Sprintf("%s.%s", u.currentExePath, ts)

	st, _ := os.Stat(u.currentExePath)
	mode := os.FileMode(0755)
	if st != nil {
		mode = st.Mode().Perm()
	}
	if err := os.WriteFile(tmp, release, mode); err != nil {
		return err
	}

	if err := os.Rename(u.currentExePath, old); err != nil {
		_ = os.Remove(tmp)
		return err
	}

	if err := os.Rename(tmp, u.currentExePath); err != nil {
		_ = os.Rename(old, u.currentExePath)
		_ = os.Remove(tmp)
		return err
	}

	u.appliedVer = version
	return nil
}
