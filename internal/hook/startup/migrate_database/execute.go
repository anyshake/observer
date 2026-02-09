package migrate_database

func (t *MigrateDatabaseStartupImpl) Execute() error {
	if err := t.ActionHandler.SchemaVersionInit(); err != nil {
		return err
	}

	// Manual database schema migrations goes here.
	// ... will be implemented in the future

	return nil
}
