package user

import (
	"errors"
	"net/http"

	"github.com/anyshake/observer/drivers/dao/tables"
	"gorm.io/gorm"
)

func (h *User) handleCreate(db *gorm.DB, currentUserId int64, req *request) (code int, msg string, err error) {
	// Check if the current user is an admin
	isAdmin, _, err := h.isAdministrator(db, currentUserId)
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
		err := errors.New("provided nonce has expired")
		return http.StatusUnauthorized, err.Error(), err
	}

	// Decrypt username with the public key and check if it is valid
	usernameBytes, err := nc.rsaKeyPair.Decrypt([]byte(req.Username), true)
	if err != nil {
		return http.StatusUnauthorized, "failed to decrypt username", err
	}
	username := string(usernameBytes)
	err = h.checkUsername(db, username)
	if err != nil {
		return http.StatusBadRequest, err.Error(), err
	}

	// Decrypt password with the public key and check if it is valid
	passwordBytes, err := nc.rsaKeyPair.Decrypt([]byte(req.Password), true)
	if err != nil {
		return http.StatusUnauthorized, "failed to decrypt password", err
	}
	password := string(passwordBytes)
	err = h.checkPassword(password)
	if err != nil {
		return http.StatusBadRequest, err.Error(), err
	}

	// Create a new user with the decrypted username and password
	adminStr := "false"
	if req.Admin {
		adminStr = "true"
	}
	sysUserModel := tables.SysUser{
		Username: username,
		Admin:    adminStr,
	}
	sysUserModel.UserId = sysUserModel.NewUserId()
	sysUserModel.Password = sysUserModel.GetHashedPassword(password)
	err = db.
		Table(sysUserModel.GetName()).
		Create(&sysUserModel).
		Error
	if err != nil {
		return http.StatusInternalServerError, "failed to create new user", err
	}

	return http.StatusOK, "successfully created a new user", nil
}
