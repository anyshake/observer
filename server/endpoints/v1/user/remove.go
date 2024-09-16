package user

import (
	"errors"
	"net/http"

	"github.com/anyshake/observer/drivers/dao/tables"
	"gorm.io/gorm"
)

func (h *User) handleRemove(db *gorm.DB, currentUserId int64, req *request) (code int, msg string, err error) {
	// Check if the current user is an admin
	isAdmin, currentUser, err := h.isAdministrator(db, currentUserId)
	if err != nil {
		return http.StatusInternalServerError, "failed to validate user privilege", err
	}
	if !isAdmin {
		err = errors.New("current action by current user is prohibited")
		return http.StatusForbidden, err.Error(), err
	}

	// Check if the nonce exists
	nc, ok := h.keyPairDataPool.GetAndDel(req.Nonce)
	if !ok {
		err := errors.New("provided nonce is invalid")
		return http.StatusUnauthorized, err.Error(), err
	}

	// Check if the nonce is still valid
	if !nc.valid() {
		err := errors.New("provided nonce is expired")
		return http.StatusUnauthorized, err.Error(), err
	}

	// Decrypt username with the public key
	usernameBytes, err := nc.rsaKeyPair.Decrypt([]byte(req.Username), true)
	if err != nil {
		return http.StatusUnauthorized, "failed to decrypt username", err
	}
	username := string(usernameBytes)

	// To keep there is at least one admin user
	if username == currentUser.Username {
		err = errors.New("removing current user is prohibited")
		return http.StatusForbidden, err.Error(), err
	}

	// Check for the existence of the selected user
	var sysUserModel tables.SysUser
	err = db.
		Table(sysUserModel.GetName()).
		Where("username = ?", username).
		First(&sysUserModel).
		Error
	if err != nil {
		return http.StatusBadRequest, "selected user does not exist", err
	}

	// Remove the selected user
	err = db.
		Table(sysUserModel.GetName()).
		Where("username = ?", username).
		Delete(sysUserModel.GetModel()).
		Error
	if err != nil {
		return http.StatusInternalServerError, "failed to remove selected user", err
	}

	return http.StatusOK, "successfully removed selected user", nil
}
