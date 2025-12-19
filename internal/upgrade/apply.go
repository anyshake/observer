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

	exePath, err := os.Executable()
	if err != nil {
		return err
	}
	exePath, err = filepath.EvalSymlinks(exePath)
	if err != nil {
		return err
	}

	// Check for write permission
	dir := filepath.Dir(exePath)
	probe, err := os.CreateTemp(dir, ".upgrade-perm-*")
	if err != nil {
		return err
	}
	probe.Close()
	_ = os.Remove(probe.Name())

	ts := time.Now().UTC().Format("20060102150405")
	tmp := fmt.Sprintf("%s.%s.new", exePath, ts)
	old := fmt.Sprintf("%s.%s", exePath, ts)

	st, _ := os.Stat(exePath)
	mode := os.FileMode(0755)
	if st != nil {
		mode = st.Mode().Perm()
	}
	if err := os.WriteFile(tmp, release, mode); err != nil {
		return err
	}

	if err := os.Rename(exePath, old); err != nil {
		_ = os.Remove(tmp)
		return err
	}

	if err := os.Rename(tmp, exePath); err != nil {
		_ = os.Rename(old, exePath)
		_ = os.Remove(tmp)
		return err
	}

	u.appliedVer = version
	return nil
}
