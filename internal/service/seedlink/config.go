package seedlink

import (
	"errors"
	"fmt"

	"github.com/anyshake/observer/config"
	"github.com/anyshake/observer/internal/dao/action"
)

type seedlinkConfigEnabledImpl struct{}

func (s *seedlinkConfigEnabledImpl) GetName() string             { return "Enable" }
func (s *seedlinkConfigEnabledImpl) GetNamespace() string        { return ID }
func (s *seedlinkConfigEnabledImpl) GetKey() string              { return "enabled" }
func (s *seedlinkConfigEnabledImpl) GetType() action.SettingType { return action.Bool }
func (s *seedlinkConfigEnabledImpl) IsRequired() bool            { return true }
func (s *seedlinkConfigEnabledImpl) GetVersion() int             { return 0 }
func (s *seedlinkConfigEnabledImpl) GetOptions() map[string]any  { return nil }
func (s *seedlinkConfigEnabledImpl) GetDefaultValue() any        { return true }
func (s *seedlinkConfigEnabledImpl) GetDescription() string {
	return "Enable SeedLink service to allow third-party client (e.g. Swarm) to connect to this station."
}
func (s *seedlinkConfigEnabledImpl) Init(handler *action.Handler) error {
	if _, err := handler.SettingsInit(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to set default SeedLink service availablity: %w", err)
	}
	return nil
}
func (s *seedlinkConfigEnabledImpl) Set(handler *action.Handler, newVal any) error {
	enabled, err := config.GetConfigValBool(newVal)
	if err != nil {
		return err
	}
	if err := handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), enabled); err != nil {
		return fmt.Errorf("failed to set SeedLink service availablity: %w", err)
	}
	return nil
}
func (s *seedlinkConfigEnabledImpl) Get(handler *action.Handler) (any, error) {
	val, _, _, err := handler.SettingsGet(s.GetNamespace(), s.GetKey())
	if err != nil {
		return nil, fmt.Errorf("failed to get SeedLink service availablity: %w", err)
	}
	enabled, ok := val.(bool)
	if !ok {
		return nil, errors.New("boolean expected")
	}
	return enabled, nil
}
func (s *seedlinkConfigEnabledImpl) Restore(handler *action.Handler) error {
	if err := handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to reset SeedLink service availablity: %w", err)
	}
	return nil
}

type seedlinkConfigCompressImpl struct{}

func (s *seedlinkConfigCompressImpl) GetName() string             { return "Use Compression" }
func (s *seedlinkConfigCompressImpl) GetNamespace() string        { return ID }
func (s *seedlinkConfigCompressImpl) GetKey() string              { return "use_compression" }
func (s *seedlinkConfigCompressImpl) GetType() action.SettingType { return action.Bool }
func (s *seedlinkConfigCompressImpl) IsRequired() bool            { return false }
func (s *seedlinkConfigCompressImpl) GetVersion() int             { return 0 }
func (s *seedlinkConfigCompressImpl) GetOptions() map[string]any  { return nil }
func (s *seedlinkConfigCompressImpl) GetDefaultValue() any        { return false }
func (s *seedlinkConfigCompressImpl) GetDescription() string {
	return "Whether to compress SeedLink data stream. If true, files will use STEIM-2 compression, encoding the count data as 30-bit instead of 32-bit. This may cause data overflow and checksum errors if a 32-bit ADC reaches full scale."
}
func (s *seedlinkConfigCompressImpl) Init(handler *action.Handler) error {
	if _, err := handler.SettingsInit(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to set default SeedLink compression flag: %w", err)
	}
	return nil
}
func (s *seedlinkConfigCompressImpl) Set(handler *action.Handler, newVal any) error {
	seedlinkCompress, err := config.GetConfigValBool(newVal)
	if err != nil {
		return err
	}
	if err := handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), seedlinkCompress); err != nil {
		return fmt.Errorf("failed to set SeedLink compression flag: %w", err)
	}
	return nil
}
func (s *seedlinkConfigCompressImpl) Get(handler *action.Handler) (any, error) {
	val, _, _, err := handler.SettingsGet(s.GetNamespace(), s.GetKey())
	if err != nil {
		return nil, fmt.Errorf("failed to get SeedLink compression flag: %w", err)
	}
	seedlinkCompress, ok := val.(bool)
	if !ok {
		return nil, errors.New("boolean expected")
	}
	return seedlinkCompress, nil
}
func (s *seedlinkConfigCompressImpl) Restore(handler *action.Handler) error {
	if err := handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to reset SeedLink compression flag: %w", err)
	}
	return nil
}

type seedlinkConfigListenHostImpl struct{}

func (s *seedlinkConfigListenHostImpl) GetName() string             { return "Listen Host" }
func (s *seedlinkConfigListenHostImpl) GetNamespace() string        { return ID }
func (s *seedlinkConfigListenHostImpl) GetKey() string              { return "listen_host" }
func (s *seedlinkConfigListenHostImpl) GetType() action.SettingType { return action.String }
func (s *seedlinkConfigListenHostImpl) IsRequired() bool            { return true }
func (s *seedlinkConfigListenHostImpl) GetVersion() int             { return 0 }
func (s *seedlinkConfigListenHostImpl) GetOptions() map[string]any  { return nil }
func (s *seedlinkConfigListenHostImpl) GetDefaultValue() any        { return "localhost" }
func (s *seedlinkConfigListenHostImpl) GetDescription() string {
	return "IP address or hostname for SeedLink server to listen, by default, the server will listen on localhost."
}
func (s *seedlinkConfigListenHostImpl) Init(handler *action.Handler) error {
	if _, err := handler.SettingsInit(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to set default SeedLink listen host: %w", err)
	}
	return nil
}
func (s *seedlinkConfigListenHostImpl) Set(handler *action.Handler, newVal any) error {
	host, err := config.GetConfigValString(newVal)
	if err != nil {
		return err
	}
	if err := handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), host); err != nil {
		return fmt.Errorf("failed to set SeedLink listen host: %w", err)
	}
	return nil
}
func (s *seedlinkConfigListenHostImpl) Get(handler *action.Handler) (any, error) {
	val, _, _, err := handler.SettingsGet(s.GetNamespace(), s.GetKey())
	if err != nil {
		return nil, fmt.Errorf("failed to get SeedLink listen host: %w", err)
	}
	host, ok := val.(string)
	if !ok {
		return nil, errors.New("string expected")
	}
	return host, nil
}
func (s *seedlinkConfigListenHostImpl) Restore(handler *action.Handler) error {
	if err := handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to reset SeedLink listen host: %w", err)
	}
	return nil
}

type seedlinkConfigListenPortImpl struct{}

func (s *seedlinkConfigListenPortImpl) GetName() string             { return "Listen Port" }
func (s *seedlinkConfigListenPortImpl) GetNamespace() string        { return ID }
func (s *seedlinkConfigListenPortImpl) GetKey() string              { return "listen_port" }
func (s *seedlinkConfigListenPortImpl) GetType() action.SettingType { return action.Int }
func (s *seedlinkConfigListenPortImpl) IsRequired() bool            { return true }
func (s *seedlinkConfigListenPortImpl) GetVersion() int             { return 0 }
func (s *seedlinkConfigListenPortImpl) GetOptions() map[string]any  { return nil }
func (s *seedlinkConfigListenPortImpl) GetDefaultValue() any        { return 18000 }
func (s *seedlinkConfigListenPortImpl) GetDescription() string {
	return "IP address or hostname for SeedLink server to listen, by default, the server will listen on localhost."
}
func (s *seedlinkConfigListenPortImpl) Init(handler *action.Handler) error {
	if _, err := handler.SettingsInit(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to set default SeedLink listen port: %w", err)
	}
	return nil
}
func (s *seedlinkConfigListenPortImpl) Set(handler *action.Handler, newVal any) error {
	port, err := config.GetConfigValInt64(newVal)
	if err != nil {
		return err
	}
	if port < 1 || port > 65535 {
		return errors.New("port must be between 1 and 65535")
	}
	if err := handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), port); err != nil {
		return fmt.Errorf("failed to set SeedLink listen port: %w", err)
	}
	return nil
}
func (s *seedlinkConfigListenPortImpl) Get(handler *action.Handler) (any, error) {
	val, _, _, err := handler.SettingsGet(s.GetNamespace(), s.GetKey())
	if err != nil {
		return nil, fmt.Errorf("failed to get SeedLink listen port: %w", err)
	}
	port, ok := val.(int64)
	if !ok {
		return nil, errors.New("integer expected")
	}
	return int(port), nil
}
func (s *seedlinkConfigListenPortImpl) Restore(handler *action.Handler) error {
	if err := handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to reset SeedLink listen port: %w", err)
	}
	return nil
}

func (s *SeedLinkServiceImpl) GetConfigConstraint() []config.IConstraint {
	return []config.IConstraint{
		&seedlinkConfigEnabledImpl{},
		&seedlinkConfigCompressImpl{},
		&seedlinkConfigListenHostImpl{},
		&seedlinkConfigListenPortImpl{},
	}
}
