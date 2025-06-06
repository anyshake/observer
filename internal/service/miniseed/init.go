package miniseed

import (
	"fmt"

	"github.com/anyshake/observer/config"
)

func (s *MiniSeedServiceImpl) Init() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, con := range s.GetConfigConstraint() {
		if err := con.Init(s.actionHandler); err != nil {
			return fmt.Errorf("failed to initlize config constraint for service %s, namespace %s, key %s: %w", ID, con.GetNamespace(), con.GetKey(), err)
		}
	}

	lifeCycle, err := (&miniSeedConfigLifeCycleImpl{}).Get(s.actionHandler)
	if err != nil {
		return err
	}
	s.lifeCycle = lifeCycle.(int)

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

	filePath, err := (&miniSeedConfigFilePathImpl{}).Get(s.actionHandler)
	if err != nil {
		return err
	}
	s.filePath = filePath.(string)

	useCompress, err := (&miniSeedConfigUseCompressImpl{}).Get(s.actionHandler)
	if err != nil {
		return err
	}
	s.useCompress = useCompress.(bool)

	return nil
}
