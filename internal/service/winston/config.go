package winston

import (
	"context"
	"errors"
	"fmt"

	"github.com/anyshake/observer/config"
	"github.com/anyshake/observer/internal/dao/action"
)

type winstonConfigEnabledImpl struct{}

func (s *winstonConfigEnabledImpl) GetName() string             { return "Enable" }
func (s *winstonConfigEnabledImpl) GetNamespace() string        { return ID }
func (s *winstonConfigEnabledImpl) GetKey() string              { return "enabled" }
func (s *winstonConfigEnabledImpl) GetType() action.SettingType { return action.Bool }
func (s *winstonConfigEnabledImpl) IsRequired() bool            { return true }
func (s *winstonConfigEnabledImpl) GetVersion() int             { return 0 }
func (s *winstonConfigEnabledImpl) GetOptions() map[string]any  { return nil }
func (s *winstonConfigEnabledImpl) GetDefaultValue() any        { return false }
func (s *winstonConfigEnabledImpl) GetDescription() string {
	return "Enable Winston service to allow third-party client (e.g. Earthworm) to connect to this station."
}
func (s *winstonConfigEnabledImpl) Init(ctx context.Context, handler *action.Handler) error {
	if _, err := handler.SettingsInit(ctx, s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to set default Winston service availability: %w", err)
	}
	return nil
}
func (s *winstonConfigEnabledImpl) Set(ctx context.Context, handler *action.Handler, newVal any) error {
	enabled, err := config.GetConfigValBool(newVal)
	if err != nil {
		return err
	}
	if err := handler.SettingsSet(ctx, s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), enabled); err != nil {
		return fmt.Errorf("failed to set Winston service availability: %w", err)
	}
	return nil
}
func (s *winstonConfigEnabledImpl) Get(ctx context.Context, handler *action.Handler) (any, error) {
	val, _, _, err := handler.SettingsGet(ctx, s.GetNamespace(), s.GetKey())
	if err != nil {
		return nil, fmt.Errorf("failed to get Winston service availability: %w", err)
	}
	enabled, ok := val.(bool)
	if !ok {
		return nil, errors.New("boolean expected")
	}
	return enabled, nil
}
func (s *winstonConfigEnabledImpl) Restore(ctx context.Context, handler *action.Handler) error {
	if err := handler.SettingsSet(ctx, s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to reset Winston service availability: %w", err)
	}
	return nil
}

type winstonConfigListenHostImpl struct{}

func (s *winstonConfigListenHostImpl) GetName() string             { return "Listen Host" }
func (s *winstonConfigListenHostImpl) GetNamespace() string        { return ID }
func (s *winstonConfigListenHostImpl) GetKey() string              { return "listen_host" }
func (s *winstonConfigListenHostImpl) GetType() action.SettingType { return action.String }
func (s *winstonConfigListenHostImpl) IsRequired() bool            { return true }
func (s *winstonConfigListenHostImpl) GetVersion() int             { return 0 }
func (s *winstonConfigListenHostImpl) GetOptions() map[string]any  { return nil }
func (s *winstonConfigListenHostImpl) GetDefaultValue() any        { return "localhost" }
func (s *winstonConfigListenHostImpl) GetDescription() string {
	return "IP address or hostname for Winston server to listen, by default, the server will listen on localhost."
}
func (s *winstonConfigListenHostImpl) Init(ctx context.Context, handler *action.Handler) error {
	if _, err := handler.SettingsInit(ctx, s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to set default Winston listen host: %w", err)
	}
	return nil
}
func (s *winstonConfigListenHostImpl) Set(ctx context.Context, handler *action.Handler, newVal any) error {
	host, err := config.GetConfigValString(newVal)
	if err != nil {
		return err
	}
	if err := handler.SettingsSet(ctx, s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), host); err != nil {
		return fmt.Errorf("failed to set Winston listen host: %w", err)
	}
	return nil
}
func (s *winstonConfigListenHostImpl) Get(ctx context.Context, handler *action.Handler) (any, error) {
	val, _, _, err := handler.SettingsGet(ctx, s.GetNamespace(), s.GetKey())
	if err != nil {
		return nil, fmt.Errorf("failed to get Winston listen host: %w", err)
	}
	host, ok := val.(string)
	if !ok {
		return nil, errors.New("string expected")
	}
	return host, nil
}
func (s *winstonConfigListenHostImpl) Restore(ctx context.Context, handler *action.Handler) error {
	if err := handler.SettingsSet(ctx, s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to reset Winston listen host: %w", err)
	}
	return nil
}

type winstonConfigListenPortImpl struct{}

func (s *winstonConfigListenPortImpl) GetName() string             { return "Listen Port" }
func (s *winstonConfigListenPortImpl) GetNamespace() string        { return ID }
func (s *winstonConfigListenPortImpl) GetKey() string              { return "listen_port" }
func (s *winstonConfigListenPortImpl) GetType() action.SettingType { return action.Int }
func (s *winstonConfigListenPortImpl) IsRequired() bool            { return true }
func (s *winstonConfigListenPortImpl) GetVersion() int             { return 0 }
func (s *winstonConfigListenPortImpl) GetOptions() map[string]any  { return nil }
func (s *winstonConfigListenPortImpl) GetDefaultValue() any        { return 16022 }
func (s *winstonConfigListenPortImpl) GetDescription() string {
	return "Port for Winston server to listen, by default, the server will listen on localhost."
}
func (s *winstonConfigListenPortImpl) Init(ctx context.Context, handler *action.Handler) error {
	if _, err := handler.SettingsInit(ctx, s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to set default Winston listen port: %w", err)
	}
	return nil
}
func (s *winstonConfigListenPortImpl) Set(ctx context.Context, handler *action.Handler, newVal any) error {
	port, err := config.GetConfigValInt64(newVal)
	if err != nil {
		return err
	}
	if port < 1 || port > 65535 {
		return errors.New("port must be between 1 and 65535")
	}
	if err := handler.SettingsSet(ctx, s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), port); err != nil {
		return fmt.Errorf("failed to set Winston listen port: %w", err)
	}
	return nil
}
func (s *winstonConfigListenPortImpl) Get(ctx context.Context, handler *action.Handler) (any, error) {
	val, _, _, err := handler.SettingsGet(ctx, s.GetNamespace(), s.GetKey())
	if err != nil {
		return nil, fmt.Errorf("failed to get Winston listen port: %w", err)
	}
	port, ok := val.(int64)
	if !ok {
		return nil, errors.New("integer expected")
	}
	return int(port), nil
}
func (s *winstonConfigListenPortImpl) Restore(ctx context.Context, handler *action.Handler) error {
	if err := handler.SettingsSet(ctx, s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to reset Winston listen port: %w", err)
	}
	return nil
}

type winstonConfigBufferSizeImpl struct{}

func (s *winstonConfigBufferSizeImpl) GetName() string             { return "Buffer Size" }
func (s *winstonConfigBufferSizeImpl) GetNamespace() string        { return ID }
func (s *winstonConfigBufferSizeImpl) GetKey() string              { return "buffer_size" }
func (s *winstonConfigBufferSizeImpl) GetType() action.SettingType { return action.Int }
func (s *winstonConfigBufferSizeImpl) IsRequired() bool            { return true }
func (s *winstonConfigBufferSizeImpl) GetVersion() int             { return 0 }
func (s *winstonConfigBufferSizeImpl) GetOptions() map[string]any  { return nil }
func (s *winstonConfigBufferSizeImpl) GetDefaultValue() any        { return 600 }
func (s *winstonConfigBufferSizeImpl) GetDescription() string {
	return "Buffer size in seconds to hold data for fast query, by default, the server will hold 10 minutes data."
}
func (s *winstonConfigBufferSizeImpl) Init(ctx context.Context, handler *action.Handler) error {
	if _, err := handler.SettingsInit(ctx, s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to set default Winston buffer size: %w", err)
	}
	return nil
}
func (s *winstonConfigBufferSizeImpl) Set(ctx context.Context, handler *action.Handler, newVal any) error {
	size, err := config.GetConfigValInt64(newVal)
	if err != nil {
		return err
	}
	if size < 1 || size > 3600 {
		return errors.New("size must be between 1 and 3600 seconds")
	}
	if err := handler.SettingsSet(ctx, s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), size); err != nil {
		return fmt.Errorf("failed to set Winston buffer size: %w", err)
	}
	return nil
}
func (s *winstonConfigBufferSizeImpl) Get(ctx context.Context, handler *action.Handler) (any, error) {
	val, _, _, err := handler.SettingsGet(ctx, s.GetNamespace(), s.GetKey())
	if err != nil {
		return nil, fmt.Errorf("failed to get Winston buffer size: %w", err)
	}
	size, ok := val.(int64)
	if !ok {
		return nil, errors.New("integer expected")
	}
	return int(size), nil
}
func (s *winstonConfigBufferSizeImpl) Restore(ctx context.Context, handler *action.Handler) error {
	if err := handler.SettingsSet(ctx, s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to reset Winston buffer size: %w", err)
	}
	return nil
}

func (s *WinstonServiceImpl) GetConfigConstraint() []config.IConstraint {
	return []config.IConstraint{
		&winstonConfigEnabledImpl{},
		&winstonConfigListenHostImpl{},
		&winstonConfigListenPortImpl{},
		&winstonConfigBufferSizeImpl{},
	}
}
