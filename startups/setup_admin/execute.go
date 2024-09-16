package setup_admin

import (
	"github.com/anyshake/observer/drivers/dao/tables"
	"github.com/anyshake/observer/startups"
	"github.com/anyshake/observer/utils/logger"
	"go.uber.org/dig"
)

const (
	DEFAULT_USERNAME = "anyshake_admin"
	DEFAULT_PASSWORD = "Anyshake@12#$"
)

func (t *SetupAdminStartupTask) Execute(depsContainer *dig.Container, options *startups.Options) error {
	if options.Config.Server.Restrict {
		var sysUserModel tables.SysUser
		options.Database.
			Table(sysUserModel.GetName()).
			First(&sysUserModel)
		hashedPassword := sysUserModel.GetHashedPassword(DEFAULT_PASSWORD)

		// Create default admin user if not exists
		if sysUserModel.UserId == 0 {
			sysUserModel.UserId = sysUserModel.NewUserId()
			sysUserModel.Username = DEFAULT_USERNAME
			sysUserModel.Password = hashedPassword
			sysUserModel.Admin = "true"
			sysUserModel.LastLogin = 0
			err := options.Database.
				Table(sysUserModel.GetName()).
				Save(&sysUserModel).
				Error
			if err != nil {
				return err
			}
			logger.GetLogger(t.GetTaskName()).Infof("created default admin user: %s, password: %s", DEFAULT_USERNAME, DEFAULT_PASSWORD)
		} else if sysUserModel.Username == DEFAULT_USERNAME {
			isDefaultPassword := sysUserModel.IsPasswordCorrect(DEFAULT_PASSWORD)
			if isDefaultPassword {
				logger.GetLogger(t.GetTaskName()).Warnf("PLEASE CHANGE DEFAULT PASSWORD FOR: %s", DEFAULT_USERNAME)
			}
		}
	}

	return nil
}
