package upgrade

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/anyshake/observer/pkg/semver"
	"golang.org/x/sys/windows"
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

	sysTime := time.Now().UTC().Format("20060102150405")
	basename := strings.TrimSuffix(exePath, filepath.Ext(exePath))
	tmp := fmt.Sprintf("%s.%s.new", basename, sysTime)
	old := fmt.Sprintf("%s.%s.exe", basename, sysTime)

	if err := os.WriteFile(tmp, release, 0600); err != nil {
		return err
	}

	// exe -> old
	if err := windows.MoveFileEx(
		windows.StringToUTF16Ptr(exePath),
		windows.StringToUTF16Ptr(old),
		windows.MOVEFILE_REPLACE_EXISTING|windows.MOVEFILE_WRITE_THROUGH,
	); err != nil {
		_ = os.Remove(tmp)
		return err
	}

	// tmp -> exe
	if err := windows.MoveFileEx(
		windows.StringToUTF16Ptr(tmp),
		windows.StringToUTF16Ptr(exePath),
		windows.MOVEFILE_REPLACE_EXISTING|windows.MOVEFILE_WRITE_THROUGH,
	); err != nil {
		_ = windows.MoveFileEx(
			windows.StringToUTF16Ptr(old),
			windows.StringToUTF16Ptr(exePath),
			windows.MOVEFILE_REPLACE_EXISTING|windows.MOVEFILE_WRITE_THROUGH,
		)
		return err
	}

	u.appliedVer = version
	return nil
}
