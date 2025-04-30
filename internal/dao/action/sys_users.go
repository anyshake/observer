package action

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/anyshake/observer/internal/dao/model"
)

func (h *Handler) SysUserHasAdmin() (bool, error) {
	if h.daoObj == nil {
		return false, errors.New("database is not opened")
	}

	var adminUsers int64
	err := h.daoObj.Database.
		Table((&model.SysUser{}).GetName(h.daoObj.GetPrefix())).
		Where("admin = ?", model.ADMIN).
		Count(&adminUsers).
		Error
	if err != nil {
		return false, fmt.Errorf("failed to check for admin users: %w", err)
	}

	return adminUsers > 0, nil
}

func (h *Handler) SysUserList() ([]model.SysUser, error) {
	if h.daoObj == nil {
		return nil, errors.New("database is not opened")
	}

	var users []model.SysUser
	err := h.daoObj.Database.
		Table((&model.SysUser{}).GetName(h.daoObj.GetPrefix())).
		Find(&users).
		Error
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	return users, nil
}

func (h *Handler) SysUserGetByUsername(username string) (model.SysUser, error) {
	if h.daoObj == nil {
		return model.SysUser{}, errors.New("database is not opened")
	}

	var user model.SysUser
	err := h.daoObj.Database.
		Table(user.GetName(h.daoObj.GetPrefix())).
		Where("username = ?", username).
		First(&user).
		Error
	if err != nil {
		return model.SysUser{}, fmt.Errorf("failed to get user info: %w", err)
	}

	return user, nil
}

func (h *Handler) SysUserGetByUserId(userId string) (model.SysUser, error) {
	if h.daoObj == nil {
		return model.SysUser{}, errors.New("database is not opened")
	}

	var user model.SysUser
	err := h.daoObj.Database.
		Table(user.GetName(h.daoObj.GetPrefix())).
		Where("user_id = ?", userId).
		First(&user).
		Error
	if err != nil {
		return model.SysUser{}, fmt.Errorf("failed to get user info: %w", err)
	}

	return user, nil
}

func (h *Handler) SysUserCreate(username, password string, isAdmin bool) (string, error) {
	if h.daoObj == nil {
		return "", errors.New("database is not opened")
	}

	if err := h.SysUserCheckUsername(username); err != nil {
		return "", err
	}
	if _, err := h.SysUserGetByUsername(username); err == nil {
		return "", fmt.Errorf("user %s already exists", username)
	}

	user := model.SysUser{
		Admin:    strconv.FormatBool(isAdmin),
		Username: username,
	}
	user.Password = user.GetHashedPassword(password)
	user.UserId = user.NewUserId()

	err := h.daoObj.Database.
		Table(user.GetName(h.daoObj.GetPrefix())).
		Save(&user).
		Error
	if err != nil {
		return "", fmt.Errorf("failed to create user: %w", err)
	}

	return user.UserId, nil
}

func (h *Handler) SysUserLogin(username, password, userAgent, userIp string) (userId string, err error) {
	if h.daoObj == nil {
		return "", errors.New("database is not opened")
	}

	user, err := h.SysUserGetByUsername(username)
	if err != nil {
		return "", fmt.Errorf("failed to get user info: %w", err)
	}

	if !user.IsPasswordCorrect(password) {
		return "", fmt.Errorf("invalid password for user %s", username)
	}

	if len(userIp) > 0 && len(userAgent) > 0 {
		user.LastLogin = time.Now().UnixMilli()
		user.UserAgent = userAgent
		user.UserIp = userIp
		err = h.SysUserUpdte(user.UserId, user)
		if err != nil {
			return "", fmt.Errorf("failed to update login status: %w", err)
		}
	}

	return user.UserId, nil
}

func (h *Handler) SysUserUpdte(userId string, user model.SysUser) error {
	if h.daoObj == nil {
		return errors.New("database is not opened")
	}

	currentUser, err := h.SysUserGetByUserId(userId)
	if err != nil {
		return err
	}

	if currentUser.Username != user.Username {
		if err := h.SysUserCheckUsername(user.Username); err != nil {
			return err
		}
	}

	currentUser.Admin = user.Admin
	currentUser.LastLogin = user.LastLogin
	currentUser.Password = user.Password
	currentUser.UserAgent = user.UserAgent
	currentUser.UserIp = user.UserIp
	currentUser.Username = user.Username
	currentUser.UpdatedAt = time.Now().UnixMilli()

	err = h.daoObj.Database.
		Table(currentUser.GetName(h.daoObj.GetPrefix())).
		Where("user_id = ?", currentUser.UserId).
		Updates(&currentUser).
		Error
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

func (h *Handler) SysUserRemove(userId string) error {
	if h.daoObj == nil {
		return errors.New("database is not opened")
	}

	var user model.SysUser
	err := h.daoObj.Database.
		Table(user.GetName(h.daoObj.GetPrefix())).
		Where("user_id = ?", userId).
		Delete(user.GetModel()).
		Error
	if err != nil {
		return fmt.Errorf("failed to remove selected user ID %s: %w", userId, err)
	}

	return nil
}

func (h *Handler) SysUserCheckUsername(username string) error {
	if len(username) < 3 || len(username) > 32 {
		return errors.New("username length should be between 3 and 32")
	}

	return nil
}

func (h *Handler) SysUserCheckPassword(password string) error {
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
