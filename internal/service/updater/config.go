package updater

import (
	"errors"
	"fmt"

	"github.com/anyshake/observer/config"
	"github.com/anyshake/observer/internal/dao/action"
	"github.com/anyshake/observer/internal/upgrade"
)

type updaterConfigEnabledImpl struct{}

func (s *updaterConfigEnabledImpl) GetName() string             { return "Enable" }
func (s *updaterConfigEnabledImpl) GetNamespace() string        { return ID }
func (s *updaterConfigEnabledImpl) GetKey() string              { return "enabled" }
func (s *updaterConfigEnabledImpl) GetType() action.SettingType { return action.Bool }
func (s *updaterConfigEnabledImpl) IsRequired() bool            { return true }
func (s *updaterConfigEnabledImpl) GetVersion() int             { return 0 }
func (s *updaterConfigEnabledImpl) GetOptions() map[string]any  { return nil }
func (s *updaterConfigEnabledImpl) GetDefaultValue() any        { return false }
func (s *updaterConfigEnabledImpl) GetDescription() string {
	return "Once enabled, the software will automatically download and apply updates in the background, and complete the upgrade upon the next startup."
}
func (s *updaterConfigEnabledImpl) Init(handler *action.Handler) error {
	if _, err := handler.SettingsInit(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to set default updater service availability: %w", err)
	}
	return nil
}
func (s *updaterConfigEnabledImpl) Set(handler *action.Handler, newVal any) error {
	enabled, err := config.GetConfigValBool(newVal)
	if err != nil {
		return err
	}
	if err := handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), enabled); err != nil {
		return fmt.Errorf("failed to set updater service availability: %w", err)
	}
	return nil
}
func (s *updaterConfigEnabledImpl) Get(handler *action.Handler) (any, error) {
	val, _, _, err := handler.SettingsGet(s.GetNamespace(), s.GetKey())
	if err != nil {
		return nil, fmt.Errorf("failed to get updater service availability: %w", err)
	}
	enabled, ok := val.(bool)
	if !ok {
		return nil, errors.New("boolean expected")
	}
	return enabled, nil
}
func (s *updaterConfigEnabledImpl) Restore(handler *action.Handler) error {
	if err := handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to reset updater service availability: %w", err)
	}
	return nil
}

type updaterConfigReleaseFetchUrlImpl struct{}

func (s *updaterConfigReleaseFetchUrlImpl) GetName() string             { return "Release Fetch URL" }
func (s *updaterConfigReleaseFetchUrlImpl) GetNamespace() string        { return ID }
func (s *updaterConfigReleaseFetchUrlImpl) GetKey() string              { return "release_fetch_url" }
func (s *updaterConfigReleaseFetchUrlImpl) GetType() action.SettingType { return action.String }
func (s *updaterConfigReleaseFetchUrlImpl) IsRequired() bool            { return true }
func (s *updaterConfigReleaseFetchUrlImpl) GetVersion() int             { return 0 }
func (s *updaterConfigReleaseFetchUrlImpl) GetOptions() map[string]any  { return nil }
func (s *updaterConfigReleaseFetchUrlImpl) GetDefaultValue() any {
	return upgrade.RELEASE_FETCH_URL_TEMPLATE
}
func (s *updaterConfigReleaseFetchUrlImpl) GetDescription() string {
	return "Specify the base URL used to fetch the latest upgrade package. By default, it points to the AnyShake GitHub Releases page."
}
func (s *updaterConfigReleaseFetchUrlImpl) Init(handler *action.Handler) error {
	if _, err := handler.SettingsInit(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to set default release package URL: %w", err)
	}
	return nil
}
func (s *updaterConfigReleaseFetchUrlImpl) Set(handler *action.Handler, newVal any) error {
	str, err := config.GetConfigValString(newVal)
	if err != nil {
		return err
	}
	if err := handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), str); err != nil {
		return fmt.Errorf("failed to set release package URL: %w", err)
	}
	return nil
}
func (s *updaterConfigReleaseFetchUrlImpl) Get(handler *action.Handler) (any, error) {
	val, _, _, err := handler.SettingsGet(s.GetNamespace(), s.GetKey())
	if err != nil {
		return nil, fmt.Errorf("failed to get release package URL: %w", err)
	}
	if _, ok := val.(string); !ok {
		return nil, errors.New("failed to assert release package URL: string expected")
	}
	return val, nil
}
func (s *updaterConfigReleaseFetchUrlImpl) Restore(handler *action.Handler) error {
	if err := handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to reset release package URL: %w", err)
	}
	return nil
}

type updaterConfigAutoRestartImpl struct{}

func (s *updaterConfigAutoRestartImpl) GetName() string             { return "Auto Restart" }
func (s *updaterConfigAutoRestartImpl) GetNamespace() string        { return ID }
func (s *updaterConfigAutoRestartImpl) GetKey() string              { return "auto_restart" }
func (s *updaterConfigAutoRestartImpl) GetType() action.SettingType { return action.Bool }
func (s *updaterConfigAutoRestartImpl) IsRequired() bool            { return true }
func (s *updaterConfigAutoRestartImpl) GetVersion() int             { return 0 }
func (s *updaterConfigAutoRestartImpl) GetOptions() map[string]any  { return nil }
func (s *updaterConfigAutoRestartImpl) GetDefaultValue() any        { return false }
func (s *updaterConfigAutoRestartImpl) GetDescription() string {
	return "Automatically restarts the entire application when upgrade is applied."
}
func (s *updaterConfigAutoRestartImpl) Init(handler *action.Handler) error {
	if _, err := handler.SettingsInit(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to set default updater auto restart availability: %w", err)
	}
	return nil
}
func (s *updaterConfigAutoRestartImpl) Set(handler *action.Handler, newVal any) error {
	enabled, err := config.GetConfigValBool(newVal)
	if err != nil {
		return err
	}
	if err := handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), enabled); err != nil {
		return fmt.Errorf("failed to set updater auto restart availability: %w", err)
	}
	return nil
}
func (s *updaterConfigAutoRestartImpl) Get(handler *action.Handler) (any, error) {
	val, _, _, err := handler.SettingsGet(s.GetNamespace(), s.GetKey())
	if err != nil {
		return nil, fmt.Errorf("failed to get updater auto restart availability: %w", err)
	}
	enabled, ok := val.(bool)
	if !ok {
		return nil, errors.New("boolean expected")
	}
	return enabled, nil
}
func (s *updaterConfigAutoRestartImpl) Restore(handler *action.Handler) error {
	if err := handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to reset updater service auto restart availability: %w", err)
	}
	return nil
}

func (s *UpdaterServiceImpl) GetConfigConstraint() []config.IConstraint {
	return []config.IConstraint{
		&updaterConfigEnabledImpl{},
		&updaterConfigReleaseFetchUrlImpl{},
		&updaterConfigAutoRestartImpl{},
	}
}
