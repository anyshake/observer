package user

import (
	"net/http"
	"strconv"

	"github.com/anyshake/observer/drivers/dao/tables"
	"gorm.io/gorm"
)

func (h *User) handleProfile(db *gorm.DB, currentUserId int64) (code int, msg string, res user, err error) {
	var sysUserModel tables.SysUser
	err = db.
		Table(sysUserModel.GetName()).
		Where("user_id = ?", currentUserId).
		First(&sysUserModel).
		Error
	if err != nil {
		return http.StatusInternalServerError, "failed to get user profile from database", res, err
	}

	res = user{
		UserId:    strconv.FormatInt(sysUserModel.UserId, 10),
		Admin:     sysUserModel.Admin == "true",
		Username:  sysUserModel.Username,
		CreatedAt: sysUserModel.CreatedAt,
		UpdatedAt: sysUserModel.UpdatedAt,
		LastLogin: sysUserModel.LastLogin,
	}
	return http.StatusOK, "successfully get user profile", res, nil
}
