package setup_admin

import (
	"fmt"
	"strings"

	"github.com/anyshake/observer/internal/dao/model"
	"github.com/anyshake/observer/pkg/logger"
)

func (t *SetupAdminStartupImpl) Execute() error {
	log := logger.GetLogger(t.GetName())

	printBanner := func(logFunc func(string, ...any), msg string) {
		border := strings.Repeat("=", len(msg)+6)
		logFunc("\n%s\n|| %s ||\n%s", border, msg, border)
	}

	hasAdmin, err := t.ActionHandler.SysUserHasAdmin()
	if err != nil {
		return err
	}

	if hasAdmin {
		if _, err := t.ActionHandler.SysUserLogin(model.DEFAULT_USERNAME, model.DEFAULT_PASSWORD, "", ""); err == nil {
			printBanner(log.Warnf, fmt.Sprintf("PLEASE CHANGE THE DEFAULT PASSWORD FOR DEFAULT USER IMMEDIATELY: %s", model.DEFAULT_USERNAME))
		}
		return nil
	}

	if _, err := t.ActionHandler.SysUserCreate(model.DEFAULT_USERNAME, model.DEFAULT_PASSWORD, true); err != nil {
		return err
	}
	printBanner(log.Infof, fmt.Sprintf("CREATED DEFAULT ADMIN USER, USERNAME: %s, PASSWORD: %s", model.DEFAULT_USERNAME, model.DEFAULT_PASSWORD))

	return nil
}
