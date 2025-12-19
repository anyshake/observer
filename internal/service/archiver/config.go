package archiver

import (
	"errors"
	"fmt"

	"github.com/anyshake/observer/config"
	"github.com/anyshake/observer/internal/dao/action"
)

type archiverConfigEnabledImpl struct{}

func (s *archiverConfigEnabledImpl) GetName() string             { return "Enable Archiver" }
func (s *archiverConfigEnabledImpl) GetNamespace() string        { return ID }
func (s *archiverConfigEnabledImpl) GetKey() string              { return "enabled" }
func (s *archiverConfigEnabledImpl) GetType() action.SettingType { return action.Bool }
func (s *archiverConfigEnabledImpl) IsRequired() bool            { return true }
func (s *archiverConfigEnabledImpl) GetVersion() int             { return 0 }
func (s *archiverConfigEnabledImpl) GetOptions() map[string]any  { return nil }
func (s *archiverConfigEnabledImpl) GetDefaultValue() any        { return true }
func (s *archiverConfigEnabledImpl) GetDescription() string {
	return "Set to true if you want to enable archiver service, otherwise the history seismic waveform query and export feature will be not available."
}
func (s *archiverConfigEnabledImpl) Init(handler *action.Handler) error {
	if _, err := handler.SettingsInit(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to set default archiver service availability: %w", err)
	}
	return nil
}
func (s *archiverConfigEnabledImpl) Set(handler *action.Handler, newVal any) error {
	enabled, err := config.GetConfigValBool(newVal)
	if err != nil {
		return err
	}
	if err := handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), enabled); err != nil {
		return fmt.Errorf("failed to set archiver service availability: %w", err)
	}
	return nil
}
func (s *archiverConfigEnabledImpl) Get(handler *action.Handler) (any, error) {
	val, _, _, err := handler.SettingsGet(s.GetNamespace(), s.GetKey())
	if err != nil {
		return nil, fmt.Errorf("failed to get archiver service availability: %w", err)
	}
	enabled, ok := val.(bool)
	if !ok {
		return nil, errors.New("boolean expected")
	}
	return enabled, nil
}
func (s *archiverConfigEnabledImpl) Restore(handler *action.Handler) error {
	if err := handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to reset archiver service availability: %w", err)
	}
	return nil
}

type archiverConfigRotationImpl struct{}

func (s *archiverConfigRotationImpl) GetName() string             { return "Data Rotation" }
func (s *archiverConfigRotationImpl) GetNamespace() string        { return ID }
func (s *archiverConfigRotationImpl) GetKey() string              { return "rotation" }
func (s *archiverConfigRotationImpl) GetType() action.SettingType { return action.Int }
func (s *archiverConfigRotationImpl) IsRequired() bool            { return true }
func (s *archiverConfigRotationImpl) GetVersion() int             { return 0 }
func (s *archiverConfigRotationImpl) GetOptions() map[string]any  { return nil }
func (s *archiverConfigRotationImpl) GetDefaultValue() any        { return 0 }
func (s *archiverConfigRotationImpl) GetDescription() string {
	return "Number of days to keep the archiver history, set to 0 to disable data rotation."
}
func (s *archiverConfigRotationImpl) Init(handler *action.Handler) error {
	if _, err := handler.SettingsInit(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to set default archiver rotation: %w", err)
	}
	return nil
}
func (s *archiverConfigRotationImpl) Set(handler *action.Handler, newVal any) error {
	archiverRotation, err := config.GetConfigValInt64(newVal)
	if err != nil {
		return err
	}
	if archiverRotation < 0 {
		return errors.New("rotation cannot be negative")
	}
	if err := handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), archiverRotation); err != nil {
		return fmt.Errorf("failed to set archiver rotation: %w", err)
	}
	return nil
}
func (s *archiverConfigRotationImpl) Get(handler *action.Handler) (any, error) {
	val, _, _, err := handler.SettingsGet(s.GetNamespace(), s.GetKey())
	if err != nil {
		return nil, fmt.Errorf("failed to get archiver rotation: %w", err)
	}
	archiverRotation, ok := val.(int64)
	if !ok {
		return nil, errors.New("integer expected")
	}
	return int(archiverRotation), nil
}
func (s *archiverConfigRotationImpl) Restore(handler *action.Handler) error {
	if err := handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to reset archiver rotation: %w", err)
	}
	return nil
}
func (s *ArchiverServiceImpl) GetConfigConstraint() []config.IConstraint {
	return []config.IConstraint{
		&archiverConfigEnabledImpl{},
		&archiverConfigRotationImpl{},
	}
}
