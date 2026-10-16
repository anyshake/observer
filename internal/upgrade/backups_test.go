package upgrade

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestExecutableBackups(t *testing.T) {
	dir := t.TempDir()
	exeName := "custom-observer"
	if runtime.GOOS == "windows" {
		exeName += ".exe"
	}
	exePath := filepath.Join(dir, exeName)
	if err := os.WriteFile(exePath, []byte("current"), 0o755); err != nil {
		t.Fatal(err)
	}
	prefix, suffix := backupNameParts(exePath)
	first := prefix + "20260102030405" + suffix
	second := prefix + "20260203040506" + suffix
	outside := filepath.Join(t.TempDir(), "outside")
	for path, data := range map[string]string{
		filepath.Join(dir, first):                          "old-1",
		filepath.Join(dir, second):                         "old-2",
		filepath.Join(dir, first+".new"):                   "temporary",
		filepath.Join(dir, prefix+"2026010203040x"+suffix): "invalid",
		outside: "outside",
	} {
		if err := os.WriteFile(path, []byte(data), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	linked := prefix + "20260304050607" + suffix
	directory := prefix + "20260405060708" + suffix
	if err := os.Mkdir(filepath.Join(dir, directory), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, filepath.Join(dir, linked)); err != nil {
		if runtime.GOOS != "windows" {
			t.Fatal(err)
		}
		linked = ""
	}
	listed, err := ListExecutableBackups(exePath)
	if err != nil || len(listed) != 2 || listed[0].Name != first || listed[1].Name != second {
		t.Fatalf("listed backups = %+v, %v", listed, err)
	}
	if listed[0].Size != 5 || listed[0].ModifiedAt <= 0 {
		t.Fatalf("backup metadata = %+v", listed[0])
	}

	for _, bad := range []string{
		"../" + first, "..\\" + first, filepath.Join(dir, first),
		exeName, first + ".new", prefix + "2026010203040x" + suffix,
	} {
		if _, err := RemoveExecutableBackups(exePath, []string{first, bad}); err == nil {
			t.Fatalf("accepted invalid backup name %q", bad)
		}
		if _, err := os.Stat(filepath.Join(dir, first)); err != nil {
			t.Fatalf("valid backup was removed before validation: %v", err)
		}
	}
	if _, err := RemoveExecutableBackups(exePath, nil); err == nil {
		t.Fatal("empty selection was accepted")
	}
	if _, err := RemoveExecutableBackups(exePath, []string{first, first}); err == nil {
		t.Fatal("duplicate backup name was accepted")
	}
	if _, err := os.Stat(filepath.Join(dir, first)); err != nil {
		t.Fatalf("valid backup was removed before duplicate validation: %v", err)
	}
	if _, err := RemoveExecutableBackups(exeName, []string{first}); err == nil {
		t.Fatal("relative executable path was accepted")
	}
	if linked != "" {
		results, err := RemoveExecutableBackups(exePath, []string{linked})
		if err != nil || len(results) != 1 || results[0].Removed || results[0].Error == "" {
			t.Fatalf("symlink removal = %+v, %v", results, err)
		}
		if _, err := os.Lstat(filepath.Join(dir, linked)); err != nil {
			t.Fatalf("symlink backup was removed: %v", err)
		}
	}
	results, err := RemoveExecutableBackups(exePath, []string{directory})
	if err != nil || len(results) != 1 || results[0].Removed || results[0].Error == "" {
		t.Fatalf("directory removal = %+v, %v", results, err)
	}
	if _, err := os.Stat(filepath.Join(dir, directory)); err != nil {
		t.Fatalf("directory backup was removed: %v", err)
	}
	results, err = RemoveExecutableBackups(exePath, []string{first, second})
	if err != nil || len(results) != 2 {
		t.Fatalf("batch removal = %+v, %v", results, err)
	}
	for i, name := range []string{first, second} {
		if results[i].Name != name || !results[i].Removed || results[i].Error != "" {
			t.Fatalf("removal %q = %+v", name, results[i])
		}
	}
	listed, err = ListExecutableBackups(exePath)
	if err != nil || len(listed) != 0 {
		t.Fatalf("remaining backups = %+v, %v", listed, err)
	}
	if content, err := os.ReadFile(outside); err != nil || string(content) != "outside" {
		t.Fatalf("outside file = %q, %v", content, err)
	}
	if content, err := os.ReadFile(exePath); err != nil || string(content) != "current" {
		t.Fatalf("current executable = %q, %v", content, err)
	}
}

func TestRemoveExecutableBackupsPartialFailure(t *testing.T) {
	exePath := filepath.Join(t.TempDir(), "observer")
	prefix, suffix := backupNameParts(exePath)
	first := prefix + "20260102030405" + suffix
	missing := prefix + "20260203040506" + suffix
	last := prefix + "20260304050607" + suffix
	for _, name := range []string{first, last} {
		if err := os.WriteFile(filepath.Join(filepath.Dir(exePath), name), []byte("backup"), 0o600); err != nil {
			t.Fatal(err)
		}
	}

	results, err := RemoveExecutableBackups(exePath, []string{first, missing, last})
	if err != nil || len(results) != 3 {
		t.Fatalf("partial removal = %+v, %v", results, err)
	}
	for i, name := range []string{first, missing, last} {
		if results[i].Name != name || results[i].Removed != (name != missing) || (results[i].Error != "") != (name == missing) {
			t.Errorf("removal %q = %+v", name, results[i])
		}
		if _, err := os.Stat(filepath.Join(filepath.Dir(exePath), name)); !os.IsNotExist(err) {
			t.Errorf("backup %q still exists: %v", name, err)
		}
	}
}
