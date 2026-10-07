package testsupport

import (
	"fmt"
	"testing"
	"time"

	"github.com/anyshake/observer/config"
	"github.com/anyshake/observer/internal/dao"
	"github.com/anyshake/observer/internal/dao/action"
	"github.com/anyshake/observer/internal/dao/model"
)

// InitStation writes the default station settings a service Init reads.
func InitStation(t *testing.T, handler *action.Handler) {
	t.Helper()

	constraints := []config.IConstraint{
		&config.StationNameConfigConstraintImpl{},
		&config.StationDescriptionConfigConstraintImpl{},
		&config.StationCountryConfigConstraintImpl{},
		&config.StationPlaceConfigConstraintImpl{},
		&config.StationAffiliationConfigConstraintImpl{},
		&config.StationStationCodeConfigConstraintImpl{},
		&config.StationNetworkCodeConfigConstraintImpl{},
		&config.StationLocationCodeConfigConstraintImpl{},
		&config.StationChannelCodesConfigConstraintImpl{},
	}
	for _, constraint := range constraints {
		if err := constraint.Init(handler); err != nil {
			t.Fatalf("init station constraint %s: %v", constraint.GetKey(), err)
		}
	}
}

// EnsureSeisTable creates the shard table that holds records for ts.
func EnsureSeisTable(t *testing.T, obj *dao.DAO, ts time.Time) {
	t.Helper()

	name := fmt.Sprintf("%sseis_records_%d", obj.GetPrefix(), ts.UTC().YearDay()%model.SEIS_RECORD_SHARDS)
	if err := obj.Database.Table(name).AutoMigrate(&model.SeisRecord{}); err != nil {
		t.Fatalf("migrate %s: %v", name, err)
	}
}
