package archiver

import (
	"fmt"

	"github.com/anyshake/observer/config"
)

func (s *ArchiverServiceImpl) Init() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, con := range s.GetConfigConstraint() {
		if err := con.Init(s.actionHandler); err != nil {
			return fmt.Errorf("failed to initlize config constraint for service %s, namespace %s, key %s: %w", ID, con.GetNamespace(), con.GetKey(), err)
		}
	}

	stationCode, err := (&config.StationStationCodeConfigConstraintImpl{}).Get(s.actionHandler)
	if err != nil {
		return err
	}
	s.stationCode = stationCode.(string)

	networkCode, err := (&config.StationNetworkCodeConfigConstraintImpl{}).Get(s.actionHandler)
	if err != nil {
		return err
	}
	s.networkCode = networkCode.(string)

	locationCode, err := (&config.StationLocationCodeConfigConstraintImpl{}).Get(s.actionHandler)
	if err != nil {
		return err
	}
	s.locationCode = locationCode.(string)

	rotation, err := (&archiverConfigRotationImpl{}).Get(s.actionHandler)
	if err != nil {
		return err
	}
	s.rotation = rotation.(int)

	return nil
}
