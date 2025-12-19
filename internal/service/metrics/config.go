package metrics

import (
	"errors"
	"fmt"

	"github.com/anyshake/observer/config"
	"github.com/anyshake/observer/internal/dao/action"
)

type metricsConfigEnabledImpl struct{}

func (s *metricsConfigEnabledImpl) GetName() string             { return "Enable" }
func (s *metricsConfigEnabledImpl) GetNamespace() string        { return ID }
func (s *metricsConfigEnabledImpl) GetKey() string              { return "enabled" }
func (s *metricsConfigEnabledImpl) GetType() action.SettingType { return action.Bool }
func (s *metricsConfigEnabledImpl) IsRequired() bool            { return true }
func (s *metricsConfigEnabledImpl) GetVersion() int             { return 0 }
func (s *metricsConfigEnabledImpl) GetOptions() map[string]any  { return nil }
func (s *metricsConfigEnabledImpl) GetDefaultValue() any        { return true }
func (s *metricsConfigEnabledImpl) GetDescription() string {
	return "The metrics service helps us improve the user experience by collecting anonymous usage data and performance metrics. This data allows us to identify issues, optimize performance, and enhance future updates. While you have the option to disable this service, keeping it enabled allows us to provide you with a better and more reliable experience."
}
func (s *metricsConfigEnabledImpl) Init(handler *action.Handler) error {
	if _, err := handler.SettingsInit(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to set default metrics service availability: %w", err)
	}
	return nil
}
func (s *metricsConfigEnabledImpl) Set(handler *action.Handler, newVal any) error {
	enabled, err := config.GetConfigValBool(newVal)
	if err != nil {
		return err
	}
	if err := handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), enabled); err != nil {
		return fmt.Errorf("failed to set SeedLink service availability: %w", err)
	}
	return nil
}
func (s *metricsConfigEnabledImpl) Get(handler *action.Handler) (any, error) {
	val, _, _, err := handler.SettingsGet(s.GetNamespace(), s.GetKey())
	if err != nil {
		return nil, fmt.Errorf("failed to get SeedLink service availability: %w", err)
	}
	enabled, ok := val.(bool)
	if !ok {
		return nil, errors.New("boolean expected")
	}
	return enabled, nil
}
func (s *metricsConfigEnabledImpl) Restore(handler *action.Handler) error {
	if err := handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to reset SeedLink service availability: %w", err)
	}
	return nil
}

func (s *MetricsServiceImpl) GetConfigConstraint() []config.IConstraint {
	return []config.IConstraint{
		&metricsConfigEnabledImpl{},
	}
}
