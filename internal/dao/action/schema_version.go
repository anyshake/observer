package action

import (
	"context"
	"time"

	"github.com/anyshake/observer/internal/dao/model"
	"gorm.io/gorm"
)

func (h *Handler) SchemaVersionInit(ctx context.Context) error {
	var schemaVersion model.SchemaVersion

	var count int64
	if err := h.daoObj.Database.WithContext(ctx).
		Table(schemaVersion.GetName(h.daoObj.GetPrefix())).
		Model(&schemaVersion).
		Count(&count).
		Error; err != nil {
		return err
	}

	if count == 0 {
		return h.daoObj.Database.WithContext(ctx).
			Table(schemaVersion.GetName(h.daoObj.GetPrefix())).
			Create(&model.SchemaVersion{
				Version:     1,
				AppliedAt:   0,
				Description: "Initial schema since v4.3.4",
			}).
			Error
	}

	return nil
}

func (h *Handler) SchemaVersionGetCurrent(ctx context.Context) (int, error) {
	var sv model.SchemaVersion
	if err := h.daoObj.Database.WithContext(ctx).
		Table(sv.GetName(h.daoObj.GetPrefix())).
		First(&sv).
		Error; err != nil {
		return 0, err
	}

	return sv.Version, nil
}

func (h *Handler) SchemaVersionUpdate(ctx context.Context, oldVer, newVer int, description string) error {
	var sv model.SchemaVersion

	result := h.daoObj.Database.WithContext(ctx).
		Table(sv.GetName(h.daoObj.GetPrefix())).
		Model(&sv).
		Where("id = ? AND version = ?", 1, oldVer).
		Updates(map[string]any{
			"version":     newVer,
			"description": description,
			"applied_at":  time.Now().UnixMilli(),
		})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrInvalidData
	}

	return nil
}
