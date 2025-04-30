package setup_admin

import (
	"fmt"
	"strings"

	"github.com/anyshake/observer/internal/dao/model"
	"github.com/anyshake/observer/pkg/logger"
)

func (t *SetupAdminStartupImpl) Execute() error {
	user, _ := t.ActionHandler.SysUserGetByUsername(model.DEFAULT_USERNAME)

	// Create default admin user if not exists
	if user.UserId == "" {
		if _, err := t.ActionHandler.SysUserCreate(model.DEFAULT_USERNAME, model.DEFAULT_PASSWORD, true); err != nil {
			return err
		}
		infoText := fmt.Sprintf("CREATED DEFAULT ADMIN USER, USERNAME: %s, PASSWORD: %s", model.DEFAULT_USERNAME, model.DEFAULT_PASSWORD)
		border := strings.Repeat("=", len(infoText)+6)
		logger.GetLogger(t.GetName()).Infof("\n%s\n|| %s ||\n%s", border, infoText, border)
	} else if user.Username == model.DEFAULT_USERNAME {
		if _, err := t.ActionHandler.SysUserLogin(model.DEFAULT_USERNAME, model.DEFAULT_PASSWORD, "", ""); err == nil {
			warnText := fmt.Sprintf("PLEASE CHANGE THE DEFAULT PASSWORD FOR DEFAULT USER IMMEDIATELY: %s", model.DEFAULT_USERNAME)
			border := strings.Repeat("=", len(warnText)+6)
			logger.GetLogger(t.GetName()).Warnf("\n%s\n|| %s ||\n%s", border, warnText, border)
		}
	}

	return nil
}
