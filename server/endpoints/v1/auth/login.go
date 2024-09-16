package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/anyshake/observer/drivers/dao/tables"
	"github.com/dchest/captcha"
	"gorm.io/gorm"
)

func (h *Auth) handleLogin(db *gorm.DB, req *request, userIp, userAgent string) (code int, userId int64, err error) {
	// Check if the nonce exists
	nc, ok := h.keyPairDataPool.GetAndDel(req.Nonce)
	if !ok {
		err := errors.New("provided nonce is invalid")
		return http.StatusUnauthorized, 0, err
	}

	// Check if the nonce is still valid
	if !nc.valid() {
		err := errors.New("provided nonce is expired")
		return http.StatusUnauthorized, 0, err
	}

	// Decode credential from base64
	credentialBytes, err := nc.rsaKeyPair.Decrypt([]byte(req.Credential), true)
	if err != nil {
		return http.StatusUnauthorized, 0, errors.New("failed to decrypt credential")
	}
	var credential map[string]any
	err = json.Unmarshal(credentialBytes, &credential)
	if err != nil {
		return http.StatusUnauthorized, 0, errors.New("failed to unserialize credential")
	}
	username, usernameOk := credential["username"].(string)
	password, passwordOk := credential["password"].(string)
	captchaId, captchaIdOk := credential["captcha_id"].(string)
	captchaVal, captchaValOk := credential["captcha_solution"].(string)
	if !usernameOk || !passwordOk || !captchaIdOk || !captchaValOk {
		err := errors.New("provided credential format is invalid")
		return http.StatusUnauthorized, 0, err
	}

	// Check if the captcha is correct
	if !captcha.VerifyString(captchaId, captchaVal) {
		err := fmt.Errorf("provided captcha solution %s is invalid", captchaVal)
		return http.StatusUnauthorized, 0, err
	}

	// Check if the credential is correct
	var sysUserModel tables.SysUser
	db.
		Table(sysUserModel.GetName()).
		Where("username = ?", username).
		First(&sysUserModel)
	if sysUserModel.UserId == 0 || !sysUserModel.IsPasswordCorrect(password) {
		err := errors.New("provided credential is invalid")
		return http.StatusUnauthorized, 0, err
	}

	// Update user last login information
	sysUserModel.LastLogin = time.Now().UnixMilli()
	sysUserModel.UserIp = userIp
	sysUserModel.UserAgent = userAgent
	err = db.
		Table(sysUserModel.GetName()).
		Save(&sysUserModel).
		Error
	if err != nil {
		return http.StatusInternalServerError, 0, errors.New("failed to update last login")
	}

	return http.StatusOK, sysUserModel.UserId, nil
}
