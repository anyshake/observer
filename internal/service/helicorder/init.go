package helicorder

import (
	"context"
	"fmt"

	"github.com/anyshake/observer/config"
)

func (s *HelicorderServiceImpl) Init(ctx context.Context) error {
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
		if err := con.Init(ctx, s.dataProvider.actionHandler); err != nil {
			return fmt.Errorf("failed to initlize config constraint for service %s, namespace %s, key %s: %w", ID, con.GetNamespace(), con.GetKey(), err)
		}
	}

	stationCode, err := (&config.StationStationCodeConfigConstraintImpl{}).Get(ctx, s.dataProvider.actionHandler)
	if err != nil {
		return err
	}
	s.dataProvider.stationCode = stationCode.(string)

	networkCode, err := (&config.StationNetworkCodeConfigConstraintImpl{}).Get(ctx, s.dataProvider.actionHandler)
	if err != nil {
		return err
	}
	s.dataProvider.networkCode = networkCode.(string)

	locationCode, err := (&config.StationLocationCodeConfigConstraintImpl{}).Get(ctx, s.dataProvider.actionHandler)
	if err != nil {
		return err
	}
	s.dataProvider.locationCode = locationCode.(string)

	filePath, err := (&helicorderConfigFilePathImpl{}).Get(ctx, s.dataProvider.actionHandler)
	if err != nil {
		return fmt.Errorf("failed to get helicorder config file path: %w", err)
	}
	s.filePath = filePath.(string)

	cacheStorage, err := (&helicorderConfigCacheStorageImpl{}).Get(ctx, s.dataProvider.actionHandler)
	if err != nil {
		return fmt.Errorf("failed to get helicorder cache storage: %w", err)
	}
	s.cacheStorage = cacheStorage.(string)

	cachePath, err := (&helicorderConfigCachePathImpl{}).Get(ctx, s.dataProvider.actionHandler)
	if err != nil {
		return fmt.Errorf("failed to get helicorder cache path: %w", err)
	}
	s.cachePath = cachePath.(string)

	queryCache, err := newQueryCache(s.cacheStorage, s.cachePath)
	if err != nil {
		return fmt.Errorf("failed to initialize helicorder cache: %w", err)
	}
	if queryCache != nil {
		if err := queryCache.Clear(); err != nil {
			return fmt.Errorf("failed to clear stale helicorder cache: %w", err)
		}
	}
	s.dataProvider.queryCache = queryCache

	imageFormat, err := (&helicorderConfigImageFormatImpl{}).Get(ctx, s.dataProvider.actionHandler)
	if err != nil {
		return fmt.Errorf("failed to get helicorder image format: %w", err)
	}
	s.imageFormat = imageFormat.(string)

	timeSpan, err := (&helicorderConfigTimeSpanImpl{}).Get(ctx, s.dataProvider.actionHandler)
	if err != nil {
		return fmt.Errorf("failed to get helicorder time span: %w", err)
	}
	s.timeSpan = timeSpan.(int)

	spanSamples, err := (&helicorderConfigSpanSamplesImpl{}).Get(ctx, s.dataProvider.actionHandler)
	if err != nil {
		return fmt.Errorf("failed to get helicorder span samples: %w", err)
	}
	s.spanSamples = spanSamples.(int)

	imageSize, err := (&helicorderConfigImageSizeImpl{}).Get(ctx, s.dataProvider.actionHandler)
	if err != nil {
		return fmt.Errorf("failed to get helicorder image size: %w", err)
	}
	s.imageSize = imageSize.(int)

	lineWidth, err := (&helicorderConfigLineWidthImpl{}).Get(ctx, s.dataProvider.actionHandler)
	if err != nil {
		return fmt.Errorf("failed to get helicorder line width: %w", err)
	}
	s.lineWidth = lineWidth.(float64)

	scaleFactors, err := (&helicorderConfigScaleFactorsImpl{}).Get(ctx, s.dataProvider.actionHandler)
	if err != nil {
		return fmt.Errorf("failed to get helicorder waveform scale factor: %w", err)
	}
	s.scaleFactors = scaleFactors.([]float64)

	return nil
}
