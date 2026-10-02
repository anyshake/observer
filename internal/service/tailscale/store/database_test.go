package store

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"path/filepath"
	"testing"
	"time"

	"github.com/anyshake/observer/internal/dao"
	"github.com/anyshake/observer/internal/dao/action"
	"github.com/anyshake/observer/internal/dao/model"
	"github.com/anyshake/observer/pkg/dbengine/engines/sqlite_modernc"
	"tailscale.com/ipn"
)

const testNamespace = "service_tailscale"

func TestDatabasePersistsState(t *testing.T) {
	databasePath := filepath.Join(t.TempDir(), "state.db")
	db, err := (&sqlite_modernc.SQLite{}).Open("", "", "", databasePath, "test_", time.Second)
	if err != nil {
		t.Fatalf("failed to open test database: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("failed to access test database: %v", err)
	}
	t.Cleanup(func() {
		_ = sqlDB.Close()
	})
	if err := db.AutoMigrate(&model.UserSettings{}); err != nil {
		t.Fatalf("failed to migrate user settings table: %v", err)
	}
	handler := action.NewHandler(&dao.DAO{Database: db})

	stateStore, err := New(handler, testNamespace)
	if err != nil {
		t.Fatalf("failed to create state store: %v", err)
	}

	key := ipn.StateKey("profile-key")
	if _, err := stateStore.ReadState(key); !errors.Is(err, ipn.ErrStateNotExist) {
		t.Fatalf("unexpected error for missing state: %v", err)
	}

	value := []byte("persisted identity")
	if err := stateStore.WriteState(key, value); err != nil {
		t.Fatalf("failed to write state: %v", err)
	}
	value[0] = 'x'

	reloadedStore, err := New(handler, testNamespace)
	if err != nil {
		t.Fatalf("failed to recreate state store: %v", err)
	}
	stored, err := reloadedStore.ReadState(key)
	if err != nil {
		t.Fatalf("failed to read persisted state: %v", err)
	}
	if !bytes.Equal(stored, []byte("persisted identity")) {
		t.Fatal("state store did not preserve the written value")
	}

	encodedState, valueType, _, err := handler.SettingsGet(testNamespace, stateKey)
	if err != nil {
		t.Fatalf("failed to read serialized state: %v", err)
	}
	if valueType != action.String {
		t.Fatalf("unexpected serialized state type: %s", valueType)
	}
	var serialized map[string]string
	if err := json.Unmarshal([]byte(encodedState.(string)), &serialized); err != nil {
		t.Fatalf("failed to decode serialized state JSON: %v", err)
	}
	if serialized[string(key)] != base64.StdEncoding.EncodeToString([]byte("persisted identity")) {
		t.Fatal("state value was not encoded as base64 in JSON")
	}

	stored[0] = 'x'
	storedAgain, err := reloadedStore.ReadState(key)
	if err != nil {
		t.Fatalf("failed to read persisted state again: %v", err)
	}
	if !bytes.Equal(storedAgain, []byte("persisted identity")) {
		t.Fatal("state store returned mutable database state")
	}

	if err := reloadedStore.WriteState(key, nil); err != nil {
		t.Fatalf("failed to delete state: %v", err)
	}
	afterDeleteStore, err := New(handler, testNamespace)
	if err != nil {
		t.Fatalf("failed to reload state after deletion: %v", err)
	}
	if _, err := afterDeleteStore.ReadState(key); !errors.Is(err, ipn.ErrStateNotExist) {
		t.Fatalf("unexpected error after deleting state: %v", err)
	}
}
