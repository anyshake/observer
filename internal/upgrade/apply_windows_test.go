package upgrade

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/anyshake/observer/pkg/semver"
)

func TestApplyUpgradeWindows(t *testing.T) {
	helper := testHelper(t, semver.New("1", "2", "0", ""))
	if err := helper.ApplyUpgrade(nil, []byte("x")); err == nil {
		t.Fatal("nil version accepted")
	}

	dir := t.TempDir()
	exe := filepath.Join(dir, "observer.exe")
	if err := os.WriteFile(exe, []byte("old"), 0o600); err != nil {
		t.Fatal(err)
	}
	helper.currentExePath = exe
	version := semver.New("1", "3", "0", "")
	if err := helper.ApplyUpgrade(version, []byte("new")); err != nil {
		t.Fatal(err)
	}
	body, err := os.ReadFile(exe)
	if err != nil || string(body) != "new" {
		t.Fatalf("replaced = %q, %v", body, err)
	}
	if err := helper.ApplyUpgrade(version, []byte("ignored")); err != nil {
		t.Fatal(err)
	}
	if err := helper.ApplyUpgrade(semver.New("1", "3", "1", ""), []byte("newer")); err != nil {
		t.Fatal(err)
	}
}
