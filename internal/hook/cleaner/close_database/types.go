package close_database

import "github.com/anyshake/observer/internal/dao"

type CloseDatabaseCleanerImpl struct {
	DAO *dao.DAO
}
