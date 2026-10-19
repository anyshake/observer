package fileutil

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestRemoveAllDoesNotFollowSymbolicLinks(t *testing.T) {
	outside := t.TempDir()
	keep := filepath.Join(outside, "keep")
	if err := os.WriteFile(keep, []byte("keep"), 0o600); err != nil {
		t.Fatal(err)
	}
	tree := filepath.Join(t.TempDir(), "tree")
	if err := os.MkdirAll(filepath.Join(tree, "child"), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, filepath.Join(tree, "link")); err != nil {
		t.Skipf("symbolic links unavailable: %v", err)
	}
	if err := RemoveAll(t.Context(), tree); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(tree); !os.IsNotExist(err) {
		t.Fatalf("tree was not removed: %v", err)
	}
	if _, err := os.Stat(keep); err != nil {
		t.Fatalf("followed the symbolic link: %v", err)
	}
}

func TestRemoveAllStopsBetweenEntries(t *testing.T) {
	tree := t.TempDir()
	for _, name := range []string{"first", "second", "third"} {
		if err := os.WriteFile(filepath.Join(tree, name), nil, 0o600); err != nil {
			t.Fatal(err)
		}
	}
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	// Cancel as soon as the traversal has removed its first entry.
	observed := &removalContext{Context: ctx, path: tree, cancel: cancel}
	if err := RemoveAll(observed, tree); !errors.Is(err, context.Canceled) {
		t.Fatalf("RemoveAll = %v; want cancellation", err)
	}
	entries, err := os.ReadDir(tree)
	if err != nil || len(entries) != 2 {
		t.Fatalf("continued deleting after cancellation: entries=%d, err=%v", len(entries), err)
	}
}

type removalContext struct {
	context.Context
	path   string
	cancel context.CancelFunc
}

func (c *removalContext) Err() error {
	entries, err := os.ReadDir(c.path)
	if err == nil && len(entries) < 3 {
		c.cancel()
	}
	return c.Context.Err()
}
