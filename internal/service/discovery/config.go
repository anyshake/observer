package mdns_discovery

import (
	"errors"
	"fmt"

	"github.com/anyshake/observer/config"
	"github.com/anyshake/observer/internal/dao/action"
	"github.com/google/uuid"
)

type discoveryConfigEnabledImpl struct{}

func (s *discoveryConfigEnabledImpl) GetName() string             { return "Enable" }
func (s *discoveryConfigEnabledImpl) GetNamespace() string        { return ID }
func (s *discoveryConfigEnabledImpl) GetKey() string              { return "enabled" }
func (s *discoveryConfigEnabledImpl) GetType() action.SettingType { return action.Bool }
func (s *discoveryConfigEnabledImpl) IsRequired() bool            { return true }
func (s *discoveryConfigEnabledImpl) GetVersion() int             { return 0 }
func (s *discoveryConfigEnabledImpl) GetOptions() map[string]any  { return nil }
func (s *discoveryConfigEnabledImpl) GetDefaultValue() any        { return true }
func (s *discoveryConfigEnabledImpl) GetDescription() string {
	return "Enable mDNS discovery service to make AnyShake Observer discoverable on the local network."
}
func (s *discoveryConfigEnabledImpl) Init(handler *action.Handler) error {
	if _, err := handler.SettingsInit(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to set default mDNS discovery service availability: %w", err)
	}
	return nil
}
func (s *discoveryConfigEnabledImpl) Set(handler *action.Handler, newVal any) error {
	enabled, err := config.GetConfigValBool(newVal)
	if err != nil {
		return err
	}
	if err := handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), enabled); err != nil {
		return fmt.Errorf("failed to set mDNS discovery service availability: %w", err)
	}
	return nil
}
func (s *discoveryConfigEnabledImpl) Get(handler *action.Handler) (any, error) {
	val, _, _, err := handler.SettingsGet(s.GetNamespace(), s.GetKey())
	if err != nil {
		return nil, fmt.Errorf("failed to get mDNS discovery service availability: %w", err)
	}
	enabled, ok := val.(bool)
	if !ok {
		return nil, errors.New("boolean expected")
	}
	return enabled, nil
}
func (s *discoveryConfigEnabledImpl) Restore(handler *action.Handler) error {
	if err := handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to reset mDNS discovery service availability: %w", err)
	}
	return nil
}

type discoveryConfigInstanceNameImpl struct{}

func (s *discoveryConfigInstanceNameImpl) GetName() string             { return "Instance Name" }
func (s *discoveryConfigInstanceNameImpl) GetNamespace() string        { return ID }
func (s *discoveryConfigInstanceNameImpl) GetKey() string              { return "instance_name" }
func (s *discoveryConfigInstanceNameImpl) GetType() action.SettingType { return action.String }
func (s *discoveryConfigInstanceNameImpl) IsRequired() bool            { return true }
func (s *discoveryConfigInstanceNameImpl) GetVersion() int             { return 0 }
func (s *discoveryConfigInstanceNameImpl) GetOptions() map[string]any  { return nil }
func (s *discoveryConfigInstanceNameImpl) GetDefaultValue() any {
	id := uuid.New().String()
	return fmt.Sprintf("anyshake-observer-%s", id[:8])
}
func (s *discoveryConfigInstanceNameImpl) GetDescription() string {
	return "An unique name to identify this AnyShake Observer instance on the local network."
}
func (s *discoveryConfigInstanceNameImpl) Init(handler *action.Handler) error {
	if _, err := handler.SettingsInit(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to set default mDNS discovery service instance name: %w", err)
	}
	return nil
}
func (s *discoveryConfigInstanceNameImpl) Set(handler *action.Handler, newVal any) error {
	host, err := config.GetConfigValString(newVal)
	if err != nil {
		return err
	}
	if err := handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), host); err != nil {
		return fmt.Errorf("failed to set mDNS discovery service instance name: %w", err)
	}
	return nil
}
func (s *discoveryConfigInstanceNameImpl) Get(handler *action.Handler) (any, error) {
	val, _, _, err := handler.SettingsGet(s.GetNamespace(), s.GetKey())
	if err != nil {
		return nil, fmt.Errorf("failed to get mDNS discovery service instance name: %w", err)
	}
	host, ok := val.(string)
	if !ok {
		return nil, errors.New("string expected")
	}
	return host, nil
}
func (s *discoveryConfigInstanceNameImpl) Restore(handler *action.Handler) error {
	if err := handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to reset mDNS discovery service instance name: %w", err)
	}
	return nil
}

func (s *DiscoveryServiceImpl) GetConfigConstraint() []config.IConstraint {
	return []config.IConstraint{
		&discoveryConfigEnabledImpl{},
		&discoveryConfigInstanceNameImpl{},
	}
}
