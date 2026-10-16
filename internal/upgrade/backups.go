package upgrade

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

// Keep deletion out of the upgrade's rename and rollback window.
var executableFilesMu sync.Mutex

type ExecutableBackup struct {
	Name       string
	Size       int64
	ModifiedAt int64
}

type ExecutableBackupRemoval struct {
	Name    string
	Removed bool
	Error   string
}

func backupNameParts(exePath string) (string, string) {
	name := filepath.Base(exePath)
	if runtime.GOOS == "windows" {
		return strings.TrimSuffix(name, filepath.Ext(name)) + ".", ".exe"
	}
	return name + ".", ""
}

func isBackupName(exePath, name string) bool {
	if name == "" || name != filepath.Base(name) || strings.ContainsAny(name, `/\`) {
		return false
	}
	prefix, suffix := backupNameParts(exePath)
	if !strings.HasPrefix(name, prefix) || !strings.HasSuffix(name, suffix) {
		return false
	}
	timestamp := strings.TrimSuffix(strings.TrimPrefix(name, prefix), suffix)
	if len(timestamp) != 14 {
		return false
	}
	for _, c := range timestamp {
		if c < '0' || c > '9' {
			return false
		}
	}
	_, err := time.Parse("20060102150405", timestamp)
	return err == nil
}

func backupPath(exePath, name string) (string, error) {
	if !filepath.IsAbs(exePath) || !isBackupName(exePath, name) {
		return "", fmt.Errorf("invalid executable backup name: %q", name)
	}
	return filepath.Join(filepath.Dir(exePath), name), nil
}

func ListExecutableBackups(exePath string) ([]ExecutableBackup, error) {
	if !filepath.IsAbs(exePath) {
		return nil, fmt.Errorf("invalid executable path")
	}
	entries, err := os.ReadDir(filepath.Dir(exePath))
	if err != nil {
		return nil, err
	}

	backups := make([]ExecutableBackup, 0)
	for _, entry := range entries {
		path, err := backupPath(exePath, entry.Name())
		if err != nil {
			continue
		}
		info, err := os.Lstat(path)
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			return nil, err
		}
		if !info.Mode().IsRegular() {
			continue
		}
		backups = append(backups, ExecutableBackup{
			Name: entry.Name(), Size: info.Size(), ModifiedAt: info.ModTime().UnixMilli(),
		})
	}
	return backups, nil
}

func RemoveExecutableBackups(exePath string, names []string) ([]ExecutableBackupRemoval, error) {
	executableFilesMu.Lock()
	defer executableFilesMu.Unlock()

	if len(names) == 0 {
		return nil, fmt.Errorf("no executable backups selected")
	}

	paths := make([]string, 0, len(names))
	seen := make(map[string]bool, len(names))
	for _, name := range names {
		if seen[name] {
			return nil, fmt.Errorf("duplicate executable backup name: %q", name)
		}
		seen[name] = true
		path, err := backupPath(exePath, name)
		if err != nil {
			return nil, err
		}
		paths = append(paths, path)
	}

	results := make([]ExecutableBackupRemoval, 0, len(paths))
	for i, path := range paths {
		result := ExecutableBackupRemoval{Name: names[i]}
		info, err := os.Lstat(path)
		if err == nil && !info.Mode().IsRegular() {
			err = fmt.Errorf("executable backup is not a regular file: %q", names[i])
		}
		if err == nil {
			err = os.Remove(path)
		}
		if err != nil {
			result.Error = err.Error()
		} else {
			result.Removed = true
		}
		results = append(results, result)
	}
	return results, nil
}
