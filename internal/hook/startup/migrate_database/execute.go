package migrate_database

import "context"

func (t *MigrateDatabaseStartupImpl) Execute(ctx context.Context) error {
	if err := t.ActionHandler.SchemaVersionInit(ctx); err != nil {
		return err
	}

	// Manual database schema migrations goes here.
	// ... will be implemented in the future

	return nil
}
