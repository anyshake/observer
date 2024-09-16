package user

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/anyshake/observer/drivers/dao/tables"
	"gorm.io/gorm"
)

func (h *User) handleList(db *gorm.DB, currentUserId int64) (code int, msg string, res []user, err error) {
	var sysUserModel tables.SysUser
	err = db.
		Table(sysUserModel.GetName()).
		Where("user_id = ?", currentUserId).
		First(&sysUserModel).
		Error
	if err != nil {
		return http.StatusInternalServerError, "failed to get current user's profile", nil, err
	}
	if sysUserModel.Admin != "true" {
		err = errors.New("list action by current user is prohibited")
		return http.StatusForbidden, err.Error(), nil, err
	}

	var sysUsers []tables.SysUser
	err = db.
		Table(sysUserModel.GetName()).
		Find(&sysUsers).
		Error
	if err != nil {
		return http.StatusInternalServerError, "failed to get users from database", nil, err
	}

	for _, u := range sysUsers {
		res = append(res, user{
			UserId:    strconv.FormatInt(u.UserId, 10),
			Admin:     u.Admin == "true",
			Username:  u.Username,
			CreatedAt: u.CreatedAt,
			UpdatedAt: u.UpdatedAt,
			LastLogin: u.LastLogin,
			UserIp:    u.UserIp,
			UserAgent: u.UserAgent,
		})
	}
	return http.StatusOK, "successfully get users", res, nil
}
