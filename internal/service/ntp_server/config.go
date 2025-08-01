package ntp_server

import (
	"errors"
	"fmt"

	"github.com/anyshake/observer/config"
	"github.com/anyshake/observer/internal/dao/action"
)

type ntpServerConfigEnabledImpl struct{}

func (s *ntpServerConfigEnabledImpl) GetName() string             { return "Enable" }
func (s *ntpServerConfigEnabledImpl) GetNamespace() string        { return ID }
func (s *ntpServerConfigEnabledImpl) GetKey() string              { return "enabled" }
func (s *ntpServerConfigEnabledImpl) GetType() action.SettingType { return action.Bool }
func (s *ntpServerConfigEnabledImpl) IsRequired() bool            { return true }
func (s *ntpServerConfigEnabledImpl) GetVersion() int             { return 0 }
func (s *ntpServerConfigEnabledImpl) GetOptions() map[string]any  { return nil }
func (s *ntpServerConfigEnabledImpl) GetDefaultValue() any        { return false }
func (s *ntpServerConfigEnabledImpl) GetDescription() string {
	return "Enable NTP server service to share precise time with other devices in the local network."
}
func (s *ntpServerConfigEnabledImpl) Init(handler *action.Handler) error {
	if _, err := handler.SettingsInit(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to set default NTP server service availablity: %w", err)
	}
	return nil
}
func (s *ntpServerConfigEnabledImpl) Set(handler *action.Handler, newVal any) error {
	enabled, err := config.GetConfigValBool(newVal)
	if err != nil {
		return err
	}
	if err := handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), enabled); err != nil {
		return fmt.Errorf("failed to set NTP server service availablity: %w", err)
	}
	return nil
}
func (s *ntpServerConfigEnabledImpl) Get(handler *action.Handler) (any, error) {
	val, _, _, err := handler.SettingsGet(s.GetNamespace(), s.GetKey())
	if err != nil {
		return nil, fmt.Errorf("failed to get NTP server service availablity: %w", err)
	}
	enabled, ok := val.(bool)
	if !ok {
		return nil, errors.New("boolean expected")
	}
	return enabled, nil
}
func (s *ntpServerConfigEnabledImpl) Restore(handler *action.Handler) error {
	if err := handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to reset NTP server service availablity: %w", err)
	}
	return nil
}

type ntpServerConfigListenHostImpl struct{}

func (s *ntpServerConfigListenHostImpl) GetName() string             { return "Listen Host" }
func (s *ntpServerConfigListenHostImpl) GetNamespace() string        { return ID }
func (s *ntpServerConfigListenHostImpl) GetKey() string              { return "listen_host" }
func (s *ntpServerConfigListenHostImpl) GetType() action.SettingType { return action.String }
func (s *ntpServerConfigListenHostImpl) IsRequired() bool            { return true }
func (s *ntpServerConfigListenHostImpl) GetVersion() int             { return 0 }
func (s *ntpServerConfigListenHostImpl) GetOptions() map[string]any  { return nil }
func (s *ntpServerConfigListenHostImpl) GetDefaultValue() any        { return "0.0.0.0" }
func (s *ntpServerConfigListenHostImpl) GetDescription() string {
	return "IP address or hostname for NTP server to listen, by default, the server will listen on 0.0.0.0."
}
func (s *ntpServerConfigListenHostImpl) Init(handler *action.Handler) error {
	if _, err := handler.SettingsInit(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to set default NTP server listen host: %w", err)
	}
	return nil
}
func (s *ntpServerConfigListenHostImpl) Set(handler *action.Handler, newVal any) error {
	host, err := config.GetConfigValString(newVal)
	if err != nil {
		return err
	}
	if err := handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), host); err != nil {
		return fmt.Errorf("failed to set NTP server listen host: %w", err)
	}
	return nil
}
func (s *ntpServerConfigListenHostImpl) Get(handler *action.Handler) (any, error) {
	val, _, _, err := handler.SettingsGet(s.GetNamespace(), s.GetKey())
	if err != nil {
		return nil, fmt.Errorf("failed to get NTP server listen host: %w", err)
	}
	host, ok := val.(string)
	if !ok {
		return nil, errors.New("string expected")
	}
	return host, nil
}
func (s *ntpServerConfigListenHostImpl) Restore(handler *action.Handler) error {
	if err := handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to reset NTP server listen host: %w", err)
	}
	return nil
}

type ntpServerConfigListenPortImpl struct{}

func (s *ntpServerConfigListenPortImpl) GetName() string             { return "Listen Port" }
func (s *ntpServerConfigListenPortImpl) GetNamespace() string        { return ID }
func (s *ntpServerConfigListenPortImpl) GetKey() string              { return "listen_port" }
func (s *ntpServerConfigListenPortImpl) GetType() action.SettingType { return action.Int }
func (s *ntpServerConfigListenPortImpl) IsRequired() bool            { return true }
func (s *ntpServerConfigListenPortImpl) GetVersion() int             { return 0 }
func (s *ntpServerConfigListenPortImpl) GetOptions() map[string]any  { return nil }
func (s *ntpServerConfigListenPortImpl) GetDefaultValue() any        { return 123 }
func (s *ntpServerConfigListenPortImpl) GetDescription() string {
	return "Port for NTP server to listen, by default, the server will listen on UDP 123."
}
func (s *ntpServerConfigListenPortImpl) Init(handler *action.Handler) error {
	if _, err := handler.SettingsInit(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to set default NTP server listen port: %w", err)
	}
	return nil
}
func (s *ntpServerConfigListenPortImpl) Set(handler *action.Handler, newVal any) error {
	port, err := config.GetConfigValInt64(newVal)
	if err != nil {
		return err
	}
	if port < 1 || port > 65535 {
		return errors.New("port must be between 1 and 65535")
	}
	if err := handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), port); err != nil {
		return fmt.Errorf("failed to set NTP server listen port: %w", err)
	}
	return nil
}
func (s *ntpServerConfigListenPortImpl) Get(handler *action.Handler) (any, error) {
	val, _, _, err := handler.SettingsGet(s.GetNamespace(), s.GetKey())
	if err != nil {
		return nil, fmt.Errorf("failed to get NTP server listen port: %w", err)
	}
	port, ok := val.(int64)
	if !ok {
		return nil, errors.New("integer expected")
	}
	return int(port), nil
}
func (s *ntpServerConfigListenPortImpl) Restore(handler *action.Handler) error {
	if err := handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to reset NTP server listen port: %w", err)
	}
	return nil
}

func (s *NtpServerServiceImpl) GetConfigConstraint() []config.IConstraint {
	return []config.IConstraint{
		&ntpServerConfigEnabledImpl{},
		&ntpServerConfigListenHostImpl{},
		&ntpServerConfigListenPortImpl{},
	}
}
