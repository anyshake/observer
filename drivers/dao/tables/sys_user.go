package tables

import (
	"github.com/anyshake/observer/drivers/dao"
	"github.com/bwmarrin/snowflake"
	"golang.org/x/crypto/bcrypt"
)

type SysUser struct {
	dao.BaseTable
	UserId    int64  `gorm:"column:user_id;index;not null;unique"` // Unique user ID generated using Snowflake
	Username  string `gorm:"column:username;index;not null;unique"`
	Password  string `gorm:"column:password;not null"` // Must be hashed using GetHashedPassword
	LastLogin int64  `gorm:"column:last_login"`
	UserIp    string `gorm:"column:user_ip"`
	UserAgent string `gorm:"column:user_agent"`
	Admin     string `gorm:"column:admin;index;not null;default:false"` // false or true in string
	UpdatedAt int64  `gorm:"column:update_at;autoUpdateTime:milli;<-:update"`
}

func (t SysUser) GetModel() any {
	return &SysUser{}
}

func (t SysUser) GetName() string {
	return "sys_user"
}

func (t SysUser) NewUserId() int64 {
	id, _ := snowflake.NewNode(1)
	return id.Generate().Int64()
}

func (t SysUser) GetHashedPassword(password string) string {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPassword)
}

func (t SysUser) IsPasswordCorrect(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(t.Password), []byte(password))
	return err == nil
}
