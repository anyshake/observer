package user

import (
	"errors"
	"regexp"

	"github.com/anyshake/observer/drivers/dao/tables"
	"gorm.io/gorm"
)

func (h *User) checkUsername(db *gorm.DB, username string) error {
	if len(username) < 4 || len(username) > 20 {
		return errors.New("username length should be between 4 and 20")
	}
	if ok, _ := regexp.MatchString("^[a-zA-Z][a-zA-Z0-9_]*$", username); !ok {
		return errors.New("username should start with an alphabet and only contain alphabets, numbers and underscores")
	}

	// Check if username already exists
	var sysUserModel tables.SysUser
	err := db.
		Table(sysUserModel.GetName()).
		Where("username = ?", username).
		First(&sysUserModel).
		Error
	if err == nil {
		err = errors.New("username conflicts with an existing user")
		return err
	}

	return nil
}

func (h *User) checkPassword(password string) error {
	if len(password) < 8 || len(password) > 32 {
		return errors.New("password length should be between 8 and 32")
	}

	uppercaseRe := regexp.MustCompile(`[A-Z]`)
	if !uppercaseRe.MatchString(password) {
		return errors.New("password must contain at least one uppercase letter")
	}

	specialCharRe := regexp.MustCompile(`[!@#~$%^&*()\-_=+[\]{};:'",.<>?/|\\]`)
	if !specialCharRe.MatchString(password) {
		return errors.New("password must contain at least one special character")
	}

	return nil
}

func (h *User) isAdministrator(db *gorm.DB, currentUserId int64) (bool, tables.SysUser, error) {
	var sysUserModel tables.SysUser
	err := db.
		Table(sysUserModel.GetName()).
		Where("user_id = ?", currentUserId).
		First(&sysUserModel).
		Error
	if err != nil {
		return false, tables.SysUser{}, err
	}
	if sysUserModel.Admin != "true" {
		return false, tables.SysUser{}, nil
	}
	return true, sysUserModel, nil
}
