package miniseed

import (
	"context"
	"fmt"

	"github.com/anyshake/observer/config"
)

func (s *MiniSeedServiceImpl) Init(ctx context.Context) error {
	if err := s.lifecycle.Lock(ctx); err != nil {
		return err
	}
	defer s.lifecycle.Unlock()

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.status.GetIsRunning() {
		return fmt.Errorf("service %s cannot be initialized while running", ID)
	}
	if err := ctx.Err(); err != nil {
		return err
	}

	for _, con := range s.GetConfigConstraint() {
		if err := con.Init(ctx, s.actionHandler); err != nil {
			return fmt.Errorf("failed to initlize config constraint for service %s, namespace %s, key %s: %w", ID, con.GetNamespace(), con.GetKey(), err)
		}
	}

	lifeCycle, err := (&miniSeedConfigLifeCycleImpl{}).Get(ctx, s.actionHandler)
	if err != nil {
		return err
	}
	s.lifeCycle = lifeCycle.(int)

	stationCode, err := (&config.StationStationCodeConfigConstraintImpl{}).Get(ctx, s.actionHandler)
	if err != nil {
		return err
	}
	s.stationCode = stationCode.(string)

	networkCode, err := (&config.StationNetworkCodeConfigConstraintImpl{}).Get(ctx, s.actionHandler)
	if err != nil {
		return err
	}
	s.networkCode = networkCode.(string)

	locationCode, err := (&config.StationLocationCodeConfigConstraintImpl{}).Get(ctx, s.actionHandler)
	if err != nil {
		return err
	}
	s.locationCode = locationCode.(string)

	filePath, err := (&miniSeedConfigFilePathImpl{}).Get(ctx, s.actionHandler)
	if err != nil {
		return err
	}
	s.filePath = filePath.(string)

	useCompress, err := (&miniSeedConfigUseCompressImpl{}).Get(ctx, s.actionHandler)
	if err != nil {
		return err
	}
	s.useCompress = useCompress.(bool)

	compressScheme, err := (&miniSeedConfigCompressSchemeImpl{}).Get(ctx, s.actionHandler)
	if err != nil {
		return err
	}
	s.compressScheme = compressScheme.(string)

	return nil
}
