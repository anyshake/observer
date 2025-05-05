package model

import (
	"fmt"

	"github.com/anyshake/observer/internal/dao"
	"github.com/bwmarrin/snowflake"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const (
	DEFAULT_USERNAME = "anyshake_admin"
	DEFAULT_PASSWORD = "Anyshake@12#$"
)

const (
	ADMIN     = "true"
	NON_ADMIN = "false"
)

type SysUser struct {
	dao.BaseTable
	UserId         string `gorm:"column:user_id;index;not null;unique"` // Unique user ID generated using Snowflake
	Username       string `gorm:"column:username;index;not null;unique"`
	HashedPassword string `gorm:"column:hashed_password;not null"` // Must be hashed using GetHashedPassword
	LastLogin      int64  `gorm:"column:last_login"`
	UserIp         string `gorm:"column:user_ip"`
	UserAgent      string `gorm:"column:user_agent"`
	IsAdmin        string `gorm:"column:is_admin;index;not null;default:false"` // false or true in string
	UpdatedAt      int64  `gorm:"column:update_at;autoUpdateTime:milli;<-:update"`
}

func (t *SysUser) GetModel() any {
	return &SysUser{}
}

func (t *SysUser) GetName(tablePrefix string) string {
	return fmt.Sprintf("%s%s", tablePrefix, "sys_users")
}

func (t *SysUser) UseAutoMigrate() bool {
	return true
}

func (t *SysUser) AddPlugins(dbObj *gorm.DB, tablePrefix string) ([]gorm.Plugin, error) {
	return nil, nil
}

func (t *SysUser) NewUserId() string {
	id, _ := snowflake.NewNode(1)
	return id.Generate().String()
}

func (t *SysUser) GetHashedPassword(password string) string {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPassword)
}

func (t *SysUser) IsPasswordCorrect(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(t.HashedPassword), []byte(password))
	return err == nil
}
