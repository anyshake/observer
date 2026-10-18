package model

import (
	"fmt"

	"github.com/anyshake/observer/internal/dao"
	"gorm.io/gorm"
)

type FDSNAccessKey struct {
	dao.BaseTable
	KeyId           string `gorm:"column:key_id;size:36;not null;uniqueIndex"`
	KeyHash         string `gorm:"column:key_hash;size:64;not null;uniqueIndex" json:"-"`
	KeyPrefix       string `gorm:"column:key_prefix;size:16;not null"`
	Remark          string `gorm:"column:remark;size:500;not null"`
	ExpiresAt       int64  `gorm:"column:expires_at;not null;default:0"` // 0 means no expiration.
	RevokedAt       int64  `gorm:"column:revoked_at;not null;default:0"`
	LastRequestedAt int64  `gorm:"column:last_requested_at;not null;default:0"`
	UserIp          string `gorm:"column:user_ip;size:45"`
	UserAgent       string `gorm:"column:user_agent;size:1024"`
}

func (t *FDSNAccessKey) GetModel() any {
	return &FDSNAccessKey{}
}

func (t *FDSNAccessKey) GetName(tablePrefix string) string {
	return fmt.Sprintf("%s%s", tablePrefix, "fdsn_access_keys")
}

func (t *FDSNAccessKey) UseAutoMigrate() bool {
	return true
}

func (t *FDSNAccessKey) AddPlugins(dbObj *gorm.DB, tablePrefix string) ([]gorm.Plugin, error) {
	return nil, nil
}
