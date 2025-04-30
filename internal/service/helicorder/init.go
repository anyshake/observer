package helicorder

import (
	"fmt"

	"github.com/anyshake/observer/config"
)

func (s *HelicorderServiceImpl) Init() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, con := range s.GetConfigConstraint() {
		if err := con.Init(s.dataProvider.actionHandler); err != nil {
			return fmt.Errorf("failed to initlize config constraint for service %s, namespace %s, key %s: %w", ID, con.GetNamespace(), con.GetKey(), err)
		}
	}

	stationCode, err := (&config.StationStationCodeConfigConstraintImpl{}).Get(s.dataProvider.actionHandler)
	if err != nil {
		return err
	}
	s.dataProvider.stationCode = stationCode.(string)

	networkCode, err := (&config.StationNetworkCodeConfigConstraintImpl{}).Get(s.dataProvider.actionHandler)
	if err != nil {
		return err
	}
	s.dataProvider.networkCode = networkCode.(string)

	locationCode, err := (&config.StationLocationCodeConfigConstraintImpl{}).Get(s.dataProvider.actionHandler)
	if err != nil {
		return err
	}
	s.dataProvider.locationCode = locationCode.(string)

	filePath, err := (&helicorderConfigFilePathImpl{}).Get(s.dataProvider.actionHandler)
	if err != nil {
		return fmt.Errorf("failed to get helicorder config file path: %w", err)
	}
	s.filePath = filePath.(string)

	imageFormat, err := (&helicorderConfigImageFormatImpl{}).Get(s.dataProvider.actionHandler)
	if err != nil {
		return fmt.Errorf("failed to get helicorder image format: %w", err)
	}
	s.imageFormat = imageFormat.(string)

	timeSpan, err := (&helicorderConfigTimeSpanImpl{}).Get(s.dataProvider.actionHandler)
	if err != nil {
		return fmt.Errorf("failed to get helicorder time span: %w", err)
	}
	s.timeSpan = timeSpan.(int)

	spanSamples, err := (&helicorderConfigSpanSamplesImpl{}).Get(s.dataProvider.actionHandler)
	if err != nil {
		return fmt.Errorf("failed to get helicorder span samples: %w", err)
	}
	s.spanSamples = spanSamples.(int)

	imageSize, err := (&helicorderConfigImageSizeImpl{}).Get(s.dataProvider.actionHandler)
	if err != nil {
		return fmt.Errorf("failed to get helicorder image size: %w", err)
	}
	s.imageSize = imageSize.(int)

	lineWidth, err := (&helicorderConfigLineWidthImpl{}).Get(s.dataProvider.actionHandler)
	if err != nil {
		return fmt.Errorf("failed to get helicorder line width: %w", err)
	}
	s.lineWidth = lineWidth.(float64)

	scaleFactor, err := (&helicorderConfigScaleFactorImpl{}).Get(s.dataProvider.actionHandler)
	if err != nil {
		return fmt.Errorf("failed to get helicorder waveform scale factor: %w", err)
	}
	s.scaleFactor = scaleFactor.(float64)

	return nil
}
