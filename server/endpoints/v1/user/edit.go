package user

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/anyshake/observer/drivers/dao/tables"
	"gorm.io/gorm"
)

func (h *User) handleEdit(db *gorm.DB, currentUserId int64, req *request) (code int, msg string, err error) {
	// Check if the current user is an admin
	isAdmin, currentUser, err := h.isAdministrator(db, currentUserId)
	if err != nil {
		return http.StatusInternalServerError, "failed to validate user privilege", err
	}
	if !isAdmin {
		err = errors.New("current action by current user is prohibited")
		return http.StatusForbidden, err.Error(), err
	}

	// Check if the nonce exists and is valid
	nc, ok := h.keyPairDataPool.GetAndDel(req.Nonce)
	if !ok {
		err := errors.New("provided nonce is invalid")
		return http.StatusUnauthorized, err.Error(), err
	}
	if !nc.valid() {
		err := errors.New("provided nonce is expired")
		return http.StatusUnauthorized, err.Error(), err
	}

	// Decrypt user ID with the public key
	userIdBytes, err := nc.rsaKeyPair.Decrypt([]byte(req.UserId), true)
	if err != nil {
		return http.StatusUnauthorized, "failed to decrypt user ID", err
	}
	userId, err := strconv.ParseInt(string(userIdBytes), 10, 64)
	if err != nil {
		return http.StatusBadRequest, "failed to parse user ID", err
	}

	// Decrypt new username with the public key and check if it is valid
	usernameBytes, err := nc.rsaKeyPair.Decrypt([]byte(req.Username), true)
	if err != nil {
		return http.StatusUnauthorized, "failed to decrypt username", err
	}
	username := string(usernameBytes)
	if userId == currentUser.UserId && req.Admin != (currentUser.Admin == "true") {
		err = errors.New("can not change admin status of an user itself")
		return http.StatusForbidden, err.Error(), err
	}

	// Try to decrypt password with the public key and check if it is valid
	var password string
	if len(req.Password) > 0 {
		passwordBytes, err := nc.rsaKeyPair.Decrypt([]byte(req.Password), true)
		if err != nil {
			return http.StatusUnauthorized, "failed to decrypt password", err
		}
		password = string(passwordBytes)
		err = h.checkPassword(password)
		if err != nil {
			return http.StatusBadRequest, err.Error(), err
		}
	}

	// Update user with provided data
	adminStr := "false"
	if req.Admin {
		adminStr = "true"
	}
	sysUserModel := tables.SysUser{
		Username: username,
		Admin:    adminStr,
	}
	if len(password) > 0 {
		sysUserModel.Password = sysUserModel.GetHashedPassword(password)
	}
	err = db.
		Table(sysUserModel.GetName()).
		Where("user_id = ?", userId).
		Updates(&sysUserModel).
		Error
	if err != nil {
		return http.StatusInternalServerError, "failed to update user with provided data", err
	}

	return http.StatusOK, "successfully updated user with provided data", nil
}
