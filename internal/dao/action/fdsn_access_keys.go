package action

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/anyshake/observer/internal/dao/model"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

var ErrInvalidFDSNAccessKey = errors.New("invalid, expired or revoked FDSN access key")

func (h *Handler) FDSNAccessKeyList(ctx context.Context) ([]model.FDSNAccessKey, error) {
	if h.daoObj == nil {
		return nil, errors.New("database is not opened")
	}
	keys := make([]model.FDSNAccessKey, 0)
	err := h.daoObj.Database.WithContext(ctx).
		Table((&model.FDSNAccessKey{}).GetName(h.daoObj.GetPrefix())).
		Omit("key_hash").Order("created_at DESC").Find(&keys).Error
	if err != nil {
		return nil, fmt.Errorf("failed to list FDSN access keys: %w", err)
	}
	return keys, nil
}

func (h *Handler) FDSNAccessKeyCreate(ctx context.Context, remark string, expiresAt *int64, now time.Time) (model.FDSNAccessKey, string, error) {
	if h.daoObj == nil {
		return model.FDSNAccessKey{}, "", errors.New("database is not opened")
	}
	if utf8.RuneCountInString(remark) > 500 {
		return model.FDSNAccessKey{}, "", errors.New("remark must not exceed 500 characters")
	}
	if expiresAt != nil && *expiresAt <= now.UnixMilli() {
		return model.FDSNAccessKey{}, "", errors.New("expiration must be in the future, or null for no expiration")
	}
	secret := make([]byte, 32)
	if _, err := rand.Read(secret); err != nil {
		return model.FDSNAccessKey{}, "", fmt.Errorf("failed to generate FDSN access key: %w", err)
	}
	key := "fdsn_" + base64.RawURLEncoding.EncodeToString(secret)
	hash := sha256.Sum256([]byte(key))
	record := model.FDSNAccessKey{
		KeyId: uuid.NewString(), KeyHash: hex.EncodeToString(hash[:]),
		KeyPrefix: key[:13], Remark: strings.TrimSpace(remark),
	}
	record.CreatedAt = now.UnixMilli()
	if expiresAt != nil {
		record.ExpiresAt = *expiresAt
	}
	if err := h.daoObj.Database.WithContext(ctx).
		Table(record.GetName(h.daoObj.GetPrefix())).Create(&record).Error; err != nil {
		return model.FDSNAccessKey{}, "", fmt.Errorf("failed to create FDSN access key: %w", err)
	}
	return record, key, nil
}

func (h *Handler) FDSNAccessKeyRevoke(ctx context.Context, keyIds []string, now time.Time) error {
	if h.daoObj == nil {
		return errors.New("database is not opened")
	}
	if len(keyIds) == 0 {
		return errors.New("at least one FDSN access key ID is required")
	}
	keyIds = lo.Uniq(keyIds)
	tableName := (&model.FDSNAccessKey{}).GetName(h.daoObj.GetPrefix())
	err := h.daoObj.Database.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var count int64
		if err := tx.Table(tableName).Where("key_id IN ?", keyIds).Count(&count).Error; err != nil {
			return err
		}
		if count != int64(len(keyIds)) {
			return gorm.ErrRecordNotFound
		}
		return tx.Table(tableName).Where("key_id IN ? AND revoked_at = 0", keyIds).
			Update("revoked_at", now.UnixMilli()).Error
	})
	if err != nil {
		return fmt.Errorf("failed to revoke FDSN access keys: %w", err)
	}
	return nil
}

func (h *Handler) FDSNAccessKeyAuthenticate(ctx context.Context, key, userIp, userAgent string, now time.Time) (string, error) {
	if h.daoObj == nil {
		return "", errors.New("database is not opened")
	}
	if len(key) != 48 || !strings.HasPrefix(key, "fdsn_") {
		return "", ErrInvalidFDSNAccessKey
	}
	hash := sha256.Sum256([]byte(key))
	var record model.FDSNAccessKey
	// Never cache authorization: every request observes committed revocations.
	err := h.daoObj.Database.WithContext(ctx).Table(record.GetName(h.daoObj.GetPrefix())).
		Where("key_hash = ? AND revoked_at = 0 AND (expires_at = 0 OR expires_at > ?)", hex.EncodeToString(hash[:]), now.UnixMilli()).
		First(&record).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return "", ErrInvalidFDSNAccessKey
	}
	if err != nil {
		return "", fmt.Errorf("failed to authenticate FDSN access key: %w", err)
	}
	if record.LastRequestedAt != 0 && now.UnixMilli()-record.LastRequestedAt < time.Minute.Milliseconds() {
		return record.KeyId, nil
	}

	for {
		if err := ctx.Err(); err != nil {
			return "", err
		}
		result := h.fdsnAccessGroup.DoChan(record.KeyId, func() (any, error) {
			err := h.recordFDSNAccessRequest(ctx, record.KeyId, userIp, userAgent, now)
			return ctx.Err(), err
		})
		select {
		case <-ctx.Done():
			return "", ctx.Err()
		case update := <-result:
			if err := ctx.Err(); err != nil {
				return "", err
			}
			// A canceled leader must not fail requests whose own contexts are
			// still valid. Retry in a new flight using a surviving request.
			if update.Val != nil && (errors.Is(update.Err, context.Canceled) || errors.Is(update.Err, context.DeadlineExceeded)) {
				continue
			}
			if update.Err != nil {
				return "", update.Err
			}
			return record.KeyId, nil
		}
	}
}

func (h *Handler) recordFDSNAccessRequest(ctx context.Context, keyId, userIp, userAgent string, now time.Time) error {
	// A request may reach a new flight with a snapshot from before the previous update.
	var latest model.FDSNAccessKey
	err := h.daoObj.Database.WithContext(ctx).Table(latest.GetName(h.daoObj.GetPrefix())).
		Where("key_id = ? AND revoked_at = 0 AND (expires_at = 0 OR expires_at > ?)", keyId, now.UnixMilli()).
		First(&latest).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrInvalidFDSNAccessKey
	}
	if err != nil {
		return fmt.Errorf("failed to authenticate FDSN access key: %w", err)
	}
	if latest.LastRequestedAt != 0 && now.UnixMilli()-latest.LastRequestedAt < time.Minute.Milliseconds() {
		return nil
	}
	// The conditional update also limits writes when multiple server instances share a database.
	err = h.daoObj.Database.WithContext(ctx).Table(latest.GetName(h.daoObj.GetPrefix())).
		Where("key_id = ? AND revoked_at = 0 AND (expires_at = 0 OR expires_at > ?)", keyId, now.UnixMilli()).
		Where("last_requested_at = 0 OR last_requested_at <= ?", now.Add(-time.Minute).UnixMilli()).
		Updates(map[string]any{
			"last_requested_at": now.UnixMilli(),
			"user_ip":           userIp,
			"user_agent":        string([]rune(userAgent)[:min(utf8.RuneCountInString(userAgent), 1024)]),
		}).Error
	if err != nil {
		return fmt.Errorf("failed to record FDSN access: %w", err)
	}
	return nil
}
