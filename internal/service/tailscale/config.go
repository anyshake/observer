package tailscale

import (
	crypto_rand "crypto/rand"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/anyshake/observer/config"
	"github.com/anyshake/observer/internal/dao/action"
)

type tailscaleConfigEnabledImpl struct{}

func (s *tailscaleConfigEnabledImpl) GetName() string             { return "Enable" }
func (s *tailscaleConfigEnabledImpl) GetNamespace() string        { return ID }
func (s *tailscaleConfigEnabledImpl) GetKey() string              { return "enabled" }
func (s *tailscaleConfigEnabledImpl) GetType() action.SettingType { return action.Bool }
func (s *tailscaleConfigEnabledImpl) IsRequired() bool            { return true }
func (s *tailscaleConfigEnabledImpl) GetVersion() int             { return 0 }
func (s *tailscaleConfigEnabledImpl) GetOptions() map[string]any  { return nil }
func (s *tailscaleConfigEnabledImpl) GetDefaultValue() any        { return false }
func (s *tailscaleConfigEnabledImpl) GetDescription() string {
	return "Enable Tailscale service to connect this device to a Tailscale or Headscale network."
}
func (s *tailscaleConfigEnabledImpl) Init(handler *action.Handler) error {
	if _, err := handler.SettingsInit(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to set default Tailscale service availability: %w", err)
	}
	return nil
}
func (s *tailscaleConfigEnabledImpl) Set(handler *action.Handler, newVal any) error {
	enabled, err := config.GetConfigValBool(newVal)
	if err != nil {
		return err
	}
	if err := handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), enabled); err != nil {
		return fmt.Errorf("failed to set Tailscale service availability: %w", err)
	}
	return nil
}
func (s *tailscaleConfigEnabledImpl) Get(handler *action.Handler) (any, error) {
	val, _, _, err := handler.SettingsGet(s.GetNamespace(), s.GetKey())
	if err != nil {
		return nil, fmt.Errorf("failed to get Tailscale service availability: %w", err)
	}
	enabled, ok := val.(bool)
	if !ok {
		return nil, errors.New("boolean expected")
	}
	return enabled, nil
}
func (s *tailscaleConfigEnabledImpl) Restore(handler *action.Handler) error {
	if err := handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to reset Tailscale service availability: %w", err)
	}
	return nil
}

type tailscaleConfigHostnameImpl struct{}

func (s *tailscaleConfigHostnameImpl) GetName() string             { return "Hostname" }
func (s *tailscaleConfigHostnameImpl) GetNamespace() string        { return ID }
func (s *tailscaleConfigHostnameImpl) GetKey() string              { return "hostname" }
func (s *tailscaleConfigHostnameImpl) GetType() action.SettingType { return action.String }
func (s *tailscaleConfigHostnameImpl) IsRequired() bool            { return true }
func (s *tailscaleConfigHostnameImpl) GetVersion() int             { return 0 }
func (s *tailscaleConfigHostnameImpl) GetOptions() map[string]any  { return nil }
func (s *tailscaleConfigHostnameImpl) GetDefaultValue() any {
	b := make([]byte, 4)
	if _, err := crypto_rand.Read(b); err != nil {
		return fmt.Sprintf("anyshake-observer-%x", uint32(time.Now().UnixNano()))
	}
	return fmt.Sprintf("anyshake-observer-%x", b)
}
func (s *tailscaleConfigHostnameImpl) GetDescription() string {
	return "Hostname used to identify this device in the Tailscale or Headscale network."
}
func (s *tailscaleConfigHostnameImpl) Init(handler *action.Handler) error {
	if _, err := handler.SettingsInit(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to set default Tailscale hostname: %w", err)
	}
	return nil
}
func (s *tailscaleConfigHostnameImpl) Set(handler *action.Handler, newVal any) error {
	hostname, err := config.GetConfigValString(newVal)
	if err != nil {
		return err
	}
	hostname = strings.TrimSpace(hostname)
	if hostname == "" {
		return errors.New("hostname is required")
	}
	if err := handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), hostname); err != nil {
		return fmt.Errorf("failed to set Tailscale hostname: %w", err)
	}
	return nil
}
func (s *tailscaleConfigHostnameImpl) Get(handler *action.Handler) (any, error) {
	val, _, _, err := handler.SettingsGet(s.GetNamespace(), s.GetKey())
	if err != nil {
		return nil, fmt.Errorf("failed to get Tailscale hostname: %w", err)
	}
	hostname, ok := val.(string)
	if !ok {
		return nil, errors.New("string expected")
	}
	return hostname, nil
}
func (s *tailscaleConfigHostnameImpl) Restore(handler *action.Handler) error {
	if err := handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to reset Tailscale hostname: %w", err)
	}
	return nil
}

type tailscaleConfigControlURLImpl struct{}

func (s *tailscaleConfigControlURLImpl) GetName() string             { return "Control URL" }
func (s *tailscaleConfigControlURLImpl) GetNamespace() string        { return ID }
func (s *tailscaleConfigControlURLImpl) GetKey() string              { return "control_url" }
func (s *tailscaleConfigControlURLImpl) GetType() action.SettingType { return action.String }
func (s *tailscaleConfigControlURLImpl) IsRequired() bool            { return true }
func (s *tailscaleConfigControlURLImpl) GetVersion() int             { return 0 }
func (s *tailscaleConfigControlURLImpl) GetOptions() map[string]any  { return nil }
func (s *tailscaleConfigControlURLImpl) GetDefaultValue() any {
	return "https://controlplane.tailscale.com"
}
func (s *tailscaleConfigControlURLImpl) GetDescription() string {
	return "URL of the Tailscale or Headscale coordination server."
}
func (s *tailscaleConfigControlURLImpl) Init(handler *action.Handler) error {
	if _, err := handler.SettingsInit(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to set default Tailscale control URL: %w", err)
	}
	return nil
}
func (s *tailscaleConfigControlURLImpl) Set(handler *action.Handler, newVal any) error {
	controlURL, err := config.GetConfigValString(newVal)
	if err != nil {
		return err
	}
	controlURL = strings.TrimSpace(controlURL)
	parsed, err := url.Parse(controlURL)
	if err != nil || parsed.Host == "" || (parsed.Scheme != "http" && parsed.Scheme != "https") {
		return errors.New("control URL must be a valid HTTP or HTTPS URL")
	}
	controlURL = strings.TrimRight(controlURL, "/")
	if err := handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), controlURL); err != nil {
		return fmt.Errorf("failed to set Tailscale control URL: %w", err)
	}
	return nil
}
func (s *tailscaleConfigControlURLImpl) Get(handler *action.Handler) (any, error) {
	val, _, _, err := handler.SettingsGet(s.GetNamespace(), s.GetKey())
	if err != nil {
		return nil, fmt.Errorf("failed to get Tailscale control URL: %w", err)
	}
	controlURL, ok := val.(string)
	if !ok {
		return nil, errors.New("string expected")
	}
	return controlURL, nil
}
func (s *tailscaleConfigControlURLImpl) Restore(handler *action.Handler) error {
	if err := handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to reset Tailscale control URL: %w", err)
	}
	return nil
}

type tailscaleConfigForwardPortsImpl struct{}

func (s *tailscaleConfigForwardPortsImpl) GetName() string             { return "Forward Ports" }
func (s *tailscaleConfigForwardPortsImpl) GetNamespace() string        { return ID }
func (s *tailscaleConfigForwardPortsImpl) GetKey() string              { return "forward_ports" }
func (s *tailscaleConfigForwardPortsImpl) GetType() action.SettingType { return action.StringArray }
func (s *tailscaleConfigForwardPortsImpl) IsRequired() bool            { return false }
func (s *tailscaleConfigForwardPortsImpl) GetVersion() int             { return 0 }
func (s *tailscaleConfigForwardPortsImpl) GetOptions() map[string]any  { return nil }
func (s *tailscaleConfigForwardPortsImpl) GetDefaultValue() any        { return []string{} }
func (s *tailscaleConfigForwardPortsImpl) GetDescription() string {
	return "Local ports to expose to the Tailscale or Headscale network, in port/protocol format (e.g. 80/tcp or 123/udp). The Observer web interface is always forwarded over TCP. Restart the service after changing this list."
}
func (s *tailscaleConfigForwardPortsImpl) Init(handler *action.Handler) error {
	if _, err := handler.SettingsInit(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to set default Tailscale forward ports: %w", err)
	}
	return nil
}
func (s *tailscaleConfigForwardPortsImpl) Set(handler *action.Handler, newVal any) error {
	ports, err := config.GetConfigValStringArray(newVal)
	if err != nil {
		return err
	}
	normalized, err := normalizePortForwardSpecs(ports)
	if err != nil {
		return err
	}

	if err := handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), normalized); err != nil {
		return fmt.Errorf("failed to set Tailscale forward ports: %w", err)
	}
	return nil
}
func (s *tailscaleConfigForwardPortsImpl) Get(handler *action.Handler) (any, error) {
	val, _, _, err := handler.SettingsGet(s.GetNamespace(), s.GetKey())
	if err != nil {
		return nil, fmt.Errorf("failed to get Tailscale forward ports: %w", err)
	}
	ports, ok := val.([]string)
	if !ok {
		return nil, errors.New("string array expected")
	}
	return ports, nil
}
func (s *tailscaleConfigForwardPortsImpl) Restore(handler *action.Handler) error {
	if err := handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to reset Tailscale forward ports: %w", err)
	}
	return nil
}

type tailscaleConfigAuthKeyImpl struct{}

func (s *tailscaleConfigAuthKeyImpl) GetName() string             { return "Auth Key" }
func (s *tailscaleConfigAuthKeyImpl) GetNamespace() string        { return ID }
func (s *tailscaleConfigAuthKeyImpl) GetKey() string              { return "auth_key" }
func (s *tailscaleConfigAuthKeyImpl) GetType() action.SettingType { return action.String }
func (s *tailscaleConfigAuthKeyImpl) IsRequired() bool            { return false }
func (s *tailscaleConfigAuthKeyImpl) GetVersion() int             { return 0 }
func (s *tailscaleConfigAuthKeyImpl) GetOptions() map[string]any  { return nil }
func (s *tailscaleConfigAuthKeyImpl) GetDefaultValue() any        { return "" }
func (s *tailscaleConfigAuthKeyImpl) GetDescription() string {
	return "Authentication key used to register this device with the Tailscale or Headscale coordination server. It is not used when a saved identity is available."
}
func (s *tailscaleConfigAuthKeyImpl) Init(handler *action.Handler) error {
	if _, err := handler.SettingsInit(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to set default Tailscale auth key: %w", err)
	}
	return nil
}
func (s *tailscaleConfigAuthKeyImpl) Set(handler *action.Handler, newVal any) error {
	authKey, err := config.GetConfigValString(newVal)
	if err != nil {
		return err
	}
	authKey = strings.TrimSpace(authKey)
	if err := handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), authKey); err != nil {
		return fmt.Errorf("failed to set Tailscale auth key: %w", err)
	}
	return nil
}
func (s *tailscaleConfigAuthKeyImpl) Get(handler *action.Handler) (any, error) {
	val, _, _, err := handler.SettingsGet(s.GetNamespace(), s.GetKey())
	if err != nil {
		return nil, fmt.Errorf("failed to get Tailscale auth key: %w", err)
	}
	authKey, ok := val.(string)
	if !ok {
		return nil, errors.New("string expected")
	}
	return authKey, nil
}
func (s *tailscaleConfigAuthKeyImpl) Restore(handler *action.Handler) error {
	if err := handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to reset Tailscale auth key: %w", err)
	}
	return nil
}

func (s *TailscaleServiceImpl) GetConfigConstraint() []config.IConstraint {
	return []config.IConstraint{
		&tailscaleConfigEnabledImpl{},
		&tailscaleConfigHostnameImpl{},
		&tailscaleConfigControlURLImpl{},
		&tailscaleConfigForwardPortsImpl{},
		&tailscaleConfigAuthKeyImpl{},
	}
}
