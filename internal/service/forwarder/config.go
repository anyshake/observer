package forwarder

import (
	"errors"
	"fmt"

	"github.com/anyshake/observer/config"
	"github.com/anyshake/observer/internal/dao/action"
)

type forwarderConfigEnabledImpl struct{}

func (s *forwarderConfigEnabledImpl) GetName() string             { return "Enable" }
func (s *forwarderConfigEnabledImpl) GetNamespace() string        { return ID }
func (s *forwarderConfigEnabledImpl) GetKey() string              { return "enabled" }
func (s *forwarderConfigEnabledImpl) GetType() action.SettingType { return action.Bool }
func (s *forwarderConfigEnabledImpl) IsRequired() bool            { return true }
func (s *forwarderConfigEnabledImpl) GetVersion() int             { return 0 }
func (s *forwarderConfigEnabledImpl) GetOptions() map[string]any  { return nil }
func (s *forwarderConfigEnabledImpl) GetDefaultValue() any        { return true }
func (s *forwarderConfigEnabledImpl) GetDescription() string {
	return "Enable Forwarder service to allow third-party client (e.g. Swarm) to connect to this station"
}
func (s *forwarderConfigEnabledImpl) Init(handler *action.Handler) error {
	if _, err := handler.SettingsInit(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to set default Forwarder service availability: %w", err)
	}
	return nil
}
func (s *forwarderConfigEnabledImpl) Set(handler *action.Handler, newVal any) error {
	enabled, err := config.GetConfigValBool(newVal)
	if err != nil {
		return err
	}
	if err := handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), enabled); err != nil {
		return fmt.Errorf("failed to set Forwarder service availability: %w", err)
	}
	return nil
}
func (s *forwarderConfigEnabledImpl) Get(handler *action.Handler) (any, error) {
	val, _, _, err := handler.SettingsGet(s.GetNamespace(), s.GetKey())
	if err != nil {
		return nil, fmt.Errorf("failed to get Forwarder service availability: %w", err)
	}
	enabled, ok := val.(bool)
	if !ok {
		return nil, errors.New("boolean expected")
	}
	return enabled, nil
}
func (s *forwarderConfigEnabledImpl) Restore(handler *action.Handler) error {
	if err := handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to reset Forwarder service availability: %w", err)
	}
	return nil
}

type forwarderConfigListenHostImpl struct{}

func (s *forwarderConfigListenHostImpl) GetName() string             { return "Listen Host" }
func (s *forwarderConfigListenHostImpl) GetNamespace() string        { return ID }
func (s *forwarderConfigListenHostImpl) GetKey() string              { return "listen_host" }
func (s *forwarderConfigListenHostImpl) GetType() action.SettingType { return action.String }
func (s *forwarderConfigListenHostImpl) IsRequired() bool            { return true }
func (s *forwarderConfigListenHostImpl) GetVersion() int             { return 0 }
func (s *forwarderConfigListenHostImpl) GetOptions() map[string]any  { return nil }
func (s *forwarderConfigListenHostImpl) GetDefaultValue() any        { return "localhost" }
func (s *forwarderConfigListenHostImpl) GetDescription() string {
	return "IP address or hostname for forwarder server to listen, by default, the server will listen on localhost."
}
func (s *forwarderConfigListenHostImpl) Init(handler *action.Handler) error {
	if _, err := handler.SettingsInit(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to set default forwarder listen host: %w", err)
	}
	return nil
}
func (s *forwarderConfigListenHostImpl) Set(handler *action.Handler, newVal any) error {
	host, err := config.GetConfigValString(newVal)
	if err != nil {
		return err
	}
	if err := handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), host); err != nil {
		return fmt.Errorf("failed to set forwarder listen host: %w", err)
	}
	return nil
}
func (s *forwarderConfigListenHostImpl) Get(handler *action.Handler) (any, error) {
	val, _, _, err := handler.SettingsGet(s.GetNamespace(), s.GetKey())
	if err != nil {
		return nil, fmt.Errorf("failed to get forwarder listen host: %w", err)
	}
	host, ok := val.(string)
	if !ok {
		return nil, errors.New("string expected")
	}
	return host, nil
}
func (s *forwarderConfigListenHostImpl) Restore(handler *action.Handler) error {
	if err := handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to reset forwarder listen host: %w", err)
	}
	return nil
}

type forwarderConfigListenPortImpl struct{}

func (s *forwarderConfigListenPortImpl) GetName() string             { return "Listen Port" }
func (s *forwarderConfigListenPortImpl) GetNamespace() string        { return ID }
func (s *forwarderConfigListenPortImpl) GetKey() string              { return "listen_port" }
func (s *forwarderConfigListenPortImpl) GetType() action.SettingType { return action.Int }
func (s *forwarderConfigListenPortImpl) IsRequired() bool            { return true }
func (s *forwarderConfigListenPortImpl) GetVersion() int             { return 0 }
func (s *forwarderConfigListenPortImpl) GetOptions() map[string]any  { return nil }
func (s *forwarderConfigListenPortImpl) GetDefaultValue() any        { return 30000 }
func (s *forwarderConfigListenPortImpl) GetDescription() string {
	return "IP address or hostname for forwarder server to listen, by default, the server will listen on localhost."
}
func (s *forwarderConfigListenPortImpl) Init(handler *action.Handler) error {
	if _, err := handler.SettingsInit(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to set default forwarder listen port: %w", err)
	}
	return nil
}
func (s *forwarderConfigListenPortImpl) Set(handler *action.Handler, newVal any) error {
	port, err := config.GetConfigValInt64(newVal)
	if err != nil {
		return err
	}
	if port < 1 || port > 65535 {
		return errors.New("port must be between 1 and 65535")
	}
	if err := handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), port); err != nil {
		return fmt.Errorf("failed to set forwarder listen port: %w", err)
	}
	return nil
}
func (s *forwarderConfigListenPortImpl) Get(handler *action.Handler) (any, error) {
	val, _, _, err := handler.SettingsGet(s.GetNamespace(), s.GetKey())
	if err != nil {
		return nil, fmt.Errorf("failed to get forwarder listen port: %w", err)
	}
	port, ok := val.(int64)
	if !ok {
		return nil, errors.New("integer expected")
	}
	return int(port), nil
}
func (s *forwarderConfigListenPortImpl) Restore(handler *action.Handler) error {
	if err := handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to reset forwarder listen port: %w", err)
	}
	return nil
}

func (s *ForwarderServiceImpl) GetConfigConstraint() []config.IConstraint {
	return []config.IConstraint{
		&forwarderConfigEnabledImpl{},
		&forwarderConfigListenHostImpl{},
		&forwarderConfigListenPortImpl{},
	}
}
