package action

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/anyshake/observer/internal/dao/model"
	"gorm.io/gorm"
)

func TestFDSNAccessKeyLifecycle(t *testing.T) {
	db, handler := openActionDB(t)
	if err := db.AutoMigrate(&model.FDSNAccessKey{}); err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	now := time.Date(2026, 10, 18, 0, 0, 0, 0, time.UTC)
	expiresAt := now.Add(2 * time.Minute).UnixMilli()
	record, key, err := handler.FDSNAccessKeyCreate(ctx, "client", &expiresAt, now)
	if err != nil {
		t.Fatal(err)
	}
	if record.KeyHash == key || record.CreatedAt != now.UnixMilli() {
		t.Fatal("key was stored as plaintext or creation time was not preserved")
	}
	if keyId, err := handler.FDSNAccessKeyAuthenticate(ctx, key, "192.0.2.1", "first", now); err != nil || keyId != record.KeyId {
		t.Fatalf("authenticate returned key id %q: %v", keyId, err)
	}
	if keyId, err := handler.FDSNAccessKeyAuthenticate(ctx, key, "192.0.2.2", "second", now.Add(59*time.Second)); err != nil || keyId != record.KeyId {
		t.Fatalf("authenticate within audit interval returned key id %q: %v", keyId, err)
	}
	keys, err := handler.FDSNAccessKeyList(ctx)
	if err != nil || len(keys) != 1 {
		t.Fatalf("list: %v, %v", keys, err)
	}
	if keys[0].UserIp != "192.0.2.1" || keys[0].UserAgent != "first" || keys[0].LastRequestedAt != now.UnixMilli() || keys[0].KeyHash != "" {
		t.Fatal("request audit changed within a minute or list exposed a key hash")
	}
	if _, err := handler.FDSNAccessKeyAuthenticate(ctx, key, "192.0.2.2", "second", now.Add(time.Minute)); err != nil {
		t.Fatal(err)
	}
	keys, err = handler.FDSNAccessKeyList(ctx)
	if err != nil || keys[0].LastRequestedAt != now.Add(time.Minute).UnixMilli() || keys[0].UserAgent != "second" {
		t.Fatalf("audit did not refresh after one minute: %v", err)
	}
	if err := handler.FDSNAccessKeyRevoke(ctx, []string{record.KeyId}, now.Add(time.Minute)); err != nil {
		t.Fatal(err)
	}
	for _, other := range []*Handler{handler, NewHandler(db)} {
		if keyId, err := other.FDSNAccessKeyAuthenticate(ctx, key, "192.0.2.3", "revoked", now.Add(61*time.Second)); !errors.Is(err, ErrInvalidFDSNAccessKey) || keyId != "" {
			t.Fatalf("revoked key remained valid within audit interval: %v", err)
		}
	}
	soon := now.Add(10 * time.Second).UnixMilli()
	_, expiring, err := handler.FDSNAccessKeyCreate(ctx, "short-lived", &soon, now)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := handler.FDSNAccessKeyAuthenticate(ctx, expiring, "192.0.2.1", "expiring", now); err != nil {
		t.Fatal(err)
	}
	if _, err := handler.FDSNAccessKeyAuthenticate(ctx, expiring, "192.0.2.1", "expired", now.Add(10*time.Second)); !errors.Is(err, ErrInvalidFDSNAccessKey) {
		t.Fatalf("key remained valid at its expiration within audit interval: %v", err)
	}
	if _, _, err := handler.FDSNAccessKeyCreate(ctx, "past", &soon, now.Add(time.Minute)); err == nil {
		t.Fatal("past expiration was accepted")
	}
}

func TestFDSNAccessKeyBatchRevoke(t *testing.T) {
	db, handler := openActionDB(t)
	if err := db.AutoMigrate(&model.FDSNAccessKey{}); err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	now := time.Date(2026, 10, 18, 0, 0, 0, 0, time.UTC)
	var ids, secrets []string
	for range 3 {
		record, secret, err := handler.FDSNAccessKeyCreate(ctx, "batch", nil, now)
		if err != nil {
			t.Fatal(err)
		}
		ids = append(ids, record.KeyId)
		secrets = append(secrets, secret)
	}
	if err := handler.FDSNAccessKeyRevoke(ctx, nil, now); err == nil {
		t.Fatal("empty batch was accepted")
	}
	if err := handler.FDSNAccessKeyRevoke(ctx, ids[:1], now); err != nil {
		t.Fatal(err)
	}
	if err := handler.FDSNAccessKeyRevoke(ctx, []string{ids[1], "missing"}, now); !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("batch containing a missing key returned %v", err)
	}
	if _, err := handler.FDSNAccessKeyAuthenticate(ctx, secrets[1], "192.0.2.1", "still-valid", now); err != nil {
		t.Fatalf("failed batch partially revoked a key: %v", err)
	}
	var updates atomic.Int32
	if err := db.Database.Callback().Update().After("gorm:update").Register("test:fdsn_batch_revoke", func(tx *gorm.DB) {
		if tx.Statement.Table == (&model.FDSNAccessKey{}).GetName(db.GetPrefix()) {
			updates.Add(1)
		}
	}); err != nil {
		t.Fatal(err)
	}
	revokedAt := now.Add(time.Second)
	if err := handler.FDSNAccessKeyRevoke(ctx, []string{ids[0], ids[1], ids[1]}, revokedAt); err != nil {
		t.Fatal(err)
	}
	if got := updates.Load(); got != 1 {
		t.Fatalf("batch used %d updates; want 1", got)
	}
	keys, err := handler.FDSNAccessKeyList(ctx)
	if err != nil {
		t.Fatal(err)
	}
	expected := map[string]int64{ids[0]: now.UnixMilli(), ids[1]: revokedAt.UnixMilli(), ids[2]: 0}
	for _, key := range keys {
		if key.RevokedAt != expected[key.KeyId] {
			t.Fatalf("key %s revoked at %d; want %d", key.KeyId, key.RevokedAt, expected[key.KeyId])
		}
	}
	if _, err := handler.FDSNAccessKeyAuthenticate(ctx, secrets[1], "192.0.2.1", "revoked", revokedAt); !errors.Is(err, ErrInvalidFDSNAccessKey) {
		t.Fatalf("batch-revoked key remained valid: %v", err)
	}
	if _, err := handler.FDSNAccessKeyAuthenticate(ctx, secrets[2], "192.0.2.1", "unselected", revokedAt); err != nil {
		t.Fatalf("unselected key was affected: %v", err)
	}
}

func TestFDSNAccessKeyConcurrentAudit(t *testing.T) {
	db, handler := openActionDB(t)
	if err := db.AutoMigrate(&model.FDSNAccessKey{}); err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	now := time.Date(2026, 10, 18, 0, 0, 0, 0, time.UTC)
	_, key, err := handler.FDSNAccessKeyCreate(ctx, "concurrent", nil, now)
	if err != nil {
		t.Fatal(err)
	}
	var updates atomic.Int32
	if err := db.Database.Callback().Update().After("gorm:update").Register("test:fdsn_audit", func(tx *gorm.DB) {
		if tx.Statement.Table == (&model.FDSNAccessKey{}).GetName(db.GetPrefix()) {
			updates.Add(1)
		}
	}); err != nil {
		t.Fatal(err)
	}
	for _, at := range []time.Time{now, now.Add(time.Minute)} {
		var wg sync.WaitGroup
		errs := make(chan error, 32)
		for range 32 {
			wg.Go(func() {
				_, err := handler.FDSNAccessKeyAuthenticate(ctx, key, "192.0.2.1", "parallel", at)
				errs <- err
			})
		}
		wg.Wait()
		close(errs)
		for err := range errs {
			if err != nil {
				t.Fatal(err)
			}
		}
	}
	if got := updates.Load(); got != 2 {
		t.Fatalf("audit updates = %d; want one per minute", got)
	}
}

func TestFDSNAccessKeyStaleAuditRead(t *testing.T) {
	db, handler := openActionDB(t)
	if err := db.AutoMigrate(&model.FDSNAccessKey{}); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	now := time.Date(2026, 10, 18, 0, 0, 0, 0, time.UTC)
	_, key, err := handler.FDSNAccessKeyCreate(ctx, "stale read", nil, now)
	if err != nil {
		t.Fatal(err)
	}
	tableName := (&model.FDSNAccessKey{}).GetName(db.GetPrefix())
	staleRead := make(chan struct{})
	releaseRead := make(chan struct{})
	var releaseOnce sync.Once
	release := func() { releaseOnce.Do(func() { close(releaseRead) }) }
	defer release()
	var delayed atomic.Bool
	if err := db.Database.Callback().Query().After("gorm:query").Register("test:fdsn_stale_read", func(tx *gorm.DB) {
		if tx.Statement.Table == tableName && delayed.CompareAndSwap(false, true) {
			close(staleRead)
			select {
			case <-releaseRead:
			case <-ctx.Done():
			}
		}
	}); err != nil {
		t.Fatal(err)
	}
	var updates atomic.Int32
	if err := db.Database.Callback().Update().After("gorm:update").Register("test:fdsn_stale_audit", func(tx *gorm.DB) {
		if tx.Statement.Table == tableName {
			updates.Add(1)
		}
	}); err != nil {
		t.Fatal(err)
	}
	finished := make(chan error, 1)
	go func() {
		_, err := handler.FDSNAccessKeyAuthenticate(ctx, key, "192.0.2.1", "delayed", now)
		finished <- err
	}()
	select {
	case <-staleRead:
	case <-ctx.Done():
		t.Fatal("timed out waiting for the initial audit read")
	}
	if _, err := handler.FDSNAccessKeyAuthenticate(ctx, key, "192.0.2.2", "recorded", now); err != nil {
		t.Fatal(err)
	}
	// Let the old snapshot reach singleflight after the first audit flight has finished.
	release()
	select {
	case err := <-finished:
		if err != nil {
			t.Fatal(err)
		}
	case <-ctx.Done():
		t.Fatal("timed out waiting for the delayed request")
	}
	if got := updates.Load(); got != 1 {
		t.Fatalf("audit updates = %d; want no additional update from the stale read", got)
	}
}
