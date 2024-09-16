package setup_admin

import (
	"github.com/anyshake/observer/startups"
	"go.uber.org/dig"
)

func (t *SetupAdminStartupTask) Provide(container *dig.Container, options *startups.Options) error {
	return nil
}
