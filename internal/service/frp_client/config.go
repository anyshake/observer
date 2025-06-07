package frp_client

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	crypto_rand "crypto/rand"

	"github.com/anyshake/observer/config"
	"github.com/anyshake/observer/internal/dao/action"
)

type frpClientConfigEnabledImpl struct{}

func (s *frpClientConfigEnabledImpl) GetName() string             { return "Enable" }
func (s *frpClientConfigEnabledImpl) GetNamespace() string        { return ID }
func (s *frpClientConfigEnabledImpl) GetKey() string              { return "enabled" }
func (s *frpClientConfigEnabledImpl) GetType() action.SettingType { return action.Bool }
func (s *frpClientConfigEnabledImpl) IsRequired() bool            { return true }
func (s *frpClientConfigEnabledImpl) GetVersion() int             { return 0 }
func (s *frpClientConfigEnabledImpl) GetOptions() map[string]any  { return nil }
func (s *frpClientConfigEnabledImpl) GetDefaultValue() any        { return false }
func (s *frpClientConfigEnabledImpl) GetDescription() string {
	return "Enable FRP client service to enable public access to Observer web interface via custom domain."
}
func (s *frpClientConfigEnabledImpl) Init(handler *action.Handler) error {
	if _, err := handler.SettingsInit(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to set default FRP client service availablity: %w", err)
	}
	return nil
}
func (s *frpClientConfigEnabledImpl) Set(handler *action.Handler, newVal any) error {
	enabled, err := config.GetConfigValBool(newVal)
	if err != nil {
		return err
	}
	if err := handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), enabled); err != nil {
		return fmt.Errorf("failed to set FRP client service availablity: %w", err)
	}
	return nil
}
func (s *frpClientConfigEnabledImpl) Get(handler *action.Handler) (any, error) {
	val, _, _, err := handler.SettingsGet(s.GetNamespace(), s.GetKey())
	if err != nil {
		return nil, fmt.Errorf("failed to get FRP client service availablity: %w", err)
	}
	enabled, ok := val.(bool)
	if !ok {
		return nil, errors.New("boolean expected")
	}
	return enabled, nil
}
func (s *frpClientConfigEnabledImpl) Restore(handler *action.Handler) error {
	if err := handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to reset FRP client service availablity: %w", err)
	}
	return nil
}

type frpClientConfigServerAddrImpl struct{}

func (s *frpClientConfigServerAddrImpl) GetName() string             { return "Server Address" }
func (s *frpClientConfigServerAddrImpl) GetNamespace() string        { return ID }
func (s *frpClientConfigServerAddrImpl) GetKey() string              { return "server_addr" }
func (s *frpClientConfigServerAddrImpl) GetType() action.SettingType { return action.String }
func (s *frpClientConfigServerAddrImpl) IsRequired() bool            { return true }
func (s *frpClientConfigServerAddrImpl) GetVersion() int             { return 0 }
func (s *frpClientConfigServerAddrImpl) GetOptions() map[string]any  { return nil }
func (s *frpClientConfigServerAddrImpl) GetDefaultValue() any        { return "example.com" }
func (s *frpClientConfigServerAddrImpl) GetDescription() string {
	return "The address of the FRP server to connect to."
}
func (s *frpClientConfigServerAddrImpl) Init(handler *action.Handler) error {
	if _, err := handler.SettingsInit(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to set default FRP client server address: %w", err)
	}
	return nil
}
func (s *frpClientConfigServerAddrImpl) Set(handler *action.Handler, newVal any) error {
	addr, err := config.GetConfigValString(newVal)
	if err != nil {
		return err
	}
	if addr == "" {
		return errors.New("server address cannot be empty")
	}
	if err := handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), addr); err != nil {
		return fmt.Errorf("failed to set server address: %w", err)
	}
	return nil
}
func (s *frpClientConfigServerAddrImpl) Get(handler *action.Handler) (any, error) {
	val, _, _, err := handler.SettingsGet(s.GetNamespace(), s.GetKey())
	if err != nil {
		return nil, fmt.Errorf("failed to get server address: %w", err)
	}
	addr, ok := val.(string)
	if !ok {
		return nil, errors.New("string expected")
	}
	return addr, nil
}
func (s *frpClientConfigServerAddrImpl) Restore(handler *action.Handler) error {
	if err := handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to reset server address: %w", err)
	}
	return nil
}

type frpClientConfigServerPortImpl struct{}

func (s *frpClientConfigServerPortImpl) GetName() string             { return "Server Port" }
func (s *frpClientConfigServerPortImpl) GetNamespace() string        { return ID }
func (s *frpClientConfigServerPortImpl) GetKey() string              { return "server_port" }
func (s *frpClientConfigServerPortImpl) GetType() action.SettingType { return action.Int }
func (s *frpClientConfigServerPortImpl) IsRequired() bool            { return true }
func (s *frpClientConfigServerPortImpl) GetVersion() int             { return 0 }
func (s *frpClientConfigServerPortImpl) GetOptions() map[string]any  { return nil }
func (s *frpClientConfigServerPortImpl) GetDefaultValue() any        { return 7000 }
func (s *frpClientConfigServerPortImpl) GetDescription() string {
	return "The port on which the FRP server is listening."
}
func (s *frpClientConfigServerPortImpl) Init(handler *action.Handler) error {
	if _, err := handler.SettingsInit(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to set default FRP client server port: %w", err)
	}
	return nil
}
func (s *frpClientConfigServerPortImpl) Set(handler *action.Handler, newVal any) error {
	port, err := config.GetConfigValInt64(newVal)
	if err != nil {
		return err
	}
	if port < 1 || port > 65535 {
		return errors.New("server port must be between 1 and 65535")
	}
	if err := handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), port); err != nil {
		return fmt.Errorf("failed to set server port: %w", err)
	}
	return nil
}
func (s *frpClientConfigServerPortImpl) Get(handler *action.Handler) (any, error) {
	val, _, _, err := handler.SettingsGet(s.GetNamespace(), s.GetKey())
	if err != nil {
		return nil, fmt.Errorf("failed to get server port: %w", err)
	}
	port, ok := val.(int64)
	if !ok {
		return nil, errors.New("integer expected")
	}
	return port, nil
}
func (s *frpClientConfigServerPortImpl) Restore(handler *action.Handler) error {
	err := handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue())
	if err != nil {
		return fmt.Errorf("failed to reset server port: %w", err)
	}
	return nil
}

type frpClientConfigDisableCustomTLSFirstByteImpl struct{}

func (s *frpClientConfigDisableCustomTLSFirstByteImpl) GetName() string {
	return "Disable Custom TLS First Byte"
}
func (s *frpClientConfigDisableCustomTLSFirstByteImpl) GetNamespace() string { return ID }
func (s *frpClientConfigDisableCustomTLSFirstByteImpl) GetKey() string {
	return "disable_custom_tls_first_byte"
}
func (s *frpClientConfigDisableCustomTLSFirstByteImpl) GetType() action.SettingType {
	return action.Bool
}
func (s *frpClientConfigDisableCustomTLSFirstByteImpl) IsRequired() bool           { return false }
func (s *frpClientConfigDisableCustomTLSFirstByteImpl) GetVersion() int            { return 0 }
func (s *frpClientConfigDisableCustomTLSFirstByteImpl) GetOptions() map[string]any { return nil }
func (s *frpClientConfigDisableCustomTLSFirstByteImpl) GetDefaultValue() any       { return true }
func (s *frpClientConfigDisableCustomTLSFirstByteImpl) GetDescription() string {
	return "Disable the custom first byte when using TLS."
}
func (s *frpClientConfigDisableCustomTLSFirstByteImpl) Init(handler *action.Handler) error {
	_, err := handler.SettingsInit(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue())
	if err != nil {
		return fmt.Errorf("failed to init disable_custom_tls_first_byte: %w", err)
	}
	return nil
}
func (s *frpClientConfigDisableCustomTLSFirstByteImpl) Set(handler *action.Handler, newVal any) error {
	b, ok := newVal.(bool)
	if !ok {
		return errors.New("value must be a boolean")
	}
	return handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), b)
}
func (s *frpClientConfigDisableCustomTLSFirstByteImpl) Get(handler *action.Handler) (any, error) {
	val, _, _, err := handler.SettingsGet(s.GetNamespace(), s.GetKey())
	if err != nil {
		return nil, fmt.Errorf("failed to get disable_custom_tls_first_byte: %w", err)
	}
	b, ok := val.(bool)
	if !ok {
		return nil, errors.New("bool expected")
	}
	return b, nil
}
func (s *frpClientConfigDisableCustomTLSFirstByteImpl) Restore(handler *action.Handler) error {
	return handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue())
}

type frpClientConfigUserImpl struct{}

func (s *frpClientConfigUserImpl) GetName() string             { return "User" }
func (s *frpClientConfigUserImpl) GetNamespace() string        { return ID }
func (s *frpClientConfigUserImpl) GetKey() string              { return "user" }
func (s *frpClientConfigUserImpl) GetType() action.SettingType { return action.String }
func (s *frpClientConfigUserImpl) IsRequired() bool            { return false }
func (s *frpClientConfigUserImpl) GetVersion() int             { return 0 }
func (s *frpClientConfigUserImpl) GetOptions() map[string]any  { return nil }
func (s *frpClientConfigUserImpl) GetDefaultValue() any        { return "" }
func (s *frpClientConfigUserImpl) GetDescription() string {
	return "User field as prefix in proxy name for distinguishing proxies."
}
func (s *frpClientConfigUserImpl) Init(handler *action.Handler) error {
	_, err := handler.SettingsInit(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue())
	if err != nil {
		return fmt.Errorf("failed to init user: %w", err)
	}
	return nil
}
func (s *frpClientConfigUserImpl) Set(handler *action.Handler, newVal any) error {
	user, err := config.GetConfigValString(newVal)
	if err != nil {
		return err
	}
	return handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), user)
}
func (s *frpClientConfigUserImpl) Get(handler *action.Handler) (any, error) {
	val, _, _, err := handler.SettingsGet(s.GetNamespace(), s.GetKey())
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	user, ok := val.(string)
	if !ok {
		return nil, errors.New("string expected")
	}
	return user, nil
}
func (s *frpClientConfigUserImpl) Restore(handler *action.Handler) error {
	return handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue())
}

type frpClientConfigTokenImpl struct{}

func (s *frpClientConfigTokenImpl) GetName() string             { return "Authentication Token" }
func (s *frpClientConfigTokenImpl) GetNamespace() string        { return ID }
func (s *frpClientConfigTokenImpl) GetKey() string              { return "token" }
func (s *frpClientConfigTokenImpl) GetType() action.SettingType { return action.String }
func (s *frpClientConfigTokenImpl) IsRequired() bool            { return false }
func (s *frpClientConfigTokenImpl) GetVersion() int             { return 0 }
func (s *frpClientConfigTokenImpl) GetOptions() map[string]any  { return nil }
func (s *frpClientConfigTokenImpl) GetDefaultValue() any        { return "<auth_token_here>" }
func (s *frpClientConfigTokenImpl) GetDescription() string {
	return "Authentication token used to validate client connection with the server."
}
func (s *frpClientConfigTokenImpl) Init(handler *action.Handler) error {
	if _, err := handler.SettingsInit(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to set default FRP client token: %w", err)
	}
	return nil
}
func (s *frpClientConfigTokenImpl) Set(h *action.Handler, v any) error {
	str, err := config.GetConfigValString(v)
	if err != nil {
		return errors.New("token must be a string")
	}
	return h.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), str)
}
func (s *frpClientConfigTokenImpl) Get(h *action.Handler) (any, error) {
	val, _, _, err := h.SettingsGet(s.GetNamespace(), s.GetKey())
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}
	str, ok := val.(string)
	if !ok {
		return nil, errors.New("string expected")
	}
	return str, nil
}
func (s *frpClientConfigTokenImpl) Restore(h *action.Handler) error {
	return h.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue())
}

type frpClientConfigPoolCountImpl struct{}

func (s *frpClientConfigPoolCountImpl) GetName() string             { return "Connection Pool Count" }
func (s *frpClientConfigPoolCountImpl) GetNamespace() string        { return ID }
func (s *frpClientConfigPoolCountImpl) GetKey() string              { return "pool_count" }
func (s *frpClientConfigPoolCountImpl) GetType() action.SettingType { return action.Int }
func (s *frpClientConfigPoolCountImpl) IsRequired() bool            { return false }
func (s *frpClientConfigPoolCountImpl) GetVersion() int             { return 0 }
func (s *frpClientConfigPoolCountImpl) GetOptions() map[string]any  { return nil }
func (s *frpClientConfigPoolCountImpl) GetDefaultValue() any        { return 5 }
func (s *frpClientConfigPoolCountImpl) GetDescription() string {
	return "Number of connections to keep in the pool for each proxy."
}
func (s *frpClientConfigPoolCountImpl) Init(handler *action.Handler) error {
	if _, err := handler.SettingsInit(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to set default pool_count: %w", err)
	}
	return nil
}
func (s *frpClientConfigPoolCountImpl) Set(h *action.Handler, v any) error {
	i, err := config.GetConfigValInt64(v)
	if err != nil {
		return err
	}
	if i < 0 {
		return errors.New("pool_count must be a non-negative integer")
	}
	return h.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), i)
}
func (s *frpClientConfigPoolCountImpl) Get(h *action.Handler) (any, error) {
	val, _, _, err := h.SettingsGet(s.GetNamespace(), s.GetKey())
	if err != nil {
		return nil, fmt.Errorf("failed to get pool_count: %w", err)
	}
	i, ok := val.(int64)
	if !ok {
		return nil, errors.New("intger expected")
	}
	return i, nil
}
func (s *frpClientConfigPoolCountImpl) Restore(h *action.Handler) error {
	return h.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue())
}

type frpClientConfigTcpMuxImpl struct{}

func (s *frpClientConfigTcpMuxImpl) GetName() string             { return "TCP Multiplexing" }
func (s *frpClientConfigTcpMuxImpl) GetNamespace() string        { return ID }
func (s *frpClientConfigTcpMuxImpl) GetKey() string              { return "tcp_mux" }
func (s *frpClientConfigTcpMuxImpl) GetType() action.SettingType { return action.Bool }
func (s *frpClientConfigTcpMuxImpl) IsRequired() bool            { return false }
func (s *frpClientConfigTcpMuxImpl) GetVersion() int             { return 0 }
func (s *frpClientConfigTcpMuxImpl) GetOptions() map[string]any  { return nil }
func (s *frpClientConfigTcpMuxImpl) GetDefaultValue() any        { return true }
func (s *frpClientConfigTcpMuxImpl) GetDescription() string {
	return "Enable TCP stream multiplexing for efficient connection reuse."
}
func (s *frpClientConfigTcpMuxImpl) Init(handler *action.Handler) error {
	if _, err := handler.SettingsInit(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to set default tcp_mux: %w", err)
	}
	return nil
}
func (s *frpClientConfigTcpMuxImpl) Set(h *action.Handler, v any) error {
	b, err := config.GetConfigValBool(v)
	if err != nil {
		return err
	}
	return h.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), b)
}
func (s *frpClientConfigTcpMuxImpl) Get(h *action.Handler) (any, error) {
	val, _, _, err := h.SettingsGet(s.GetNamespace(), s.GetKey())
	if err != nil {
		return nil, fmt.Errorf("failed to get tcp_mux: %w", err)
	}
	b, ok := val.(bool)
	if !ok {
		return nil, errors.New("bool expected")
	}
	return b, nil
}
func (s *frpClientConfigTcpMuxImpl) Restore(h *action.Handler) error {
	return h.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue())
}

type frpClientConfigTlsEnableImpl struct{}

func (s *frpClientConfigTlsEnableImpl) GetName() string             { return "TLS Enable" }
func (s *frpClientConfigTlsEnableImpl) GetNamespace() string        { return ID }
func (s *frpClientConfigTlsEnableImpl) GetKey() string              { return "tls_enable" }
func (s *frpClientConfigTlsEnableImpl) GetType() action.SettingType { return action.Bool }
func (s *frpClientConfigTlsEnableImpl) IsRequired() bool            { return false }
func (s *frpClientConfigTlsEnableImpl) GetVersion() int             { return 0 }
func (s *frpClientConfigTlsEnableImpl) GetOptions() map[string]any  { return nil }
func (s *frpClientConfigTlsEnableImpl) GetDefaultValue() any        { return true }
func (s *frpClientConfigTlsEnableImpl) GetDescription() string {
	return "Enable TLS encryption for control connection."
}
func (s *frpClientConfigTlsEnableImpl) Init(handler *action.Handler) error {
	if _, err := handler.SettingsInit(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to set default tls_enable: %w", err)
	}
	return nil
}
func (s *frpClientConfigTlsEnableImpl) Set(h *action.Handler, v any) error {
	b, err := config.GetConfigValBool(v)
	if err != nil {
		return err
	}
	return h.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), b)
}
func (s *frpClientConfigTlsEnableImpl) Get(h *action.Handler) (any, error) {
	val, _, _, err := h.SettingsGet(s.GetNamespace(), s.GetKey())
	if err != nil {
		return nil, fmt.Errorf("failed to get tls_enable: %w", err)
	}
	b, ok := val.(bool)
	if !ok {
		return nil, errors.New("bool expected")
	}
	return b, nil
}
func (s *frpClientConfigTlsEnableImpl) Restore(h *action.Handler) error {
	return h.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue())
}

type frpClientConfigProtocolImpl struct{}

func (s *frpClientConfigProtocolImpl) GetName() string             { return "Protocol" }
func (s *frpClientConfigProtocolImpl) GetNamespace() string        { return ID }
func (s *frpClientConfigProtocolImpl) GetKey() string              { return "protocol" }
func (s *frpClientConfigProtocolImpl) GetType() action.SettingType { return action.String }
func (s *frpClientConfigProtocolImpl) IsRequired() bool            { return true }
func (s *frpClientConfigProtocolImpl) GetVersion() int             { return 0 }
func (s *frpClientConfigProtocolImpl) GetOptions() map[string]any {
	return map[string]any{
		"TCP":       "tcp",
		"KCP":       "kcp",
		"QUIC":      "quic",
		"WebSocket": "websocket",
	}
}
func (s *frpClientConfigProtocolImpl) GetDefaultValue() any { return "tcp" }
func (s *frpClientConfigProtocolImpl) GetDescription() string {
	return "Protocol used to communicate with the server."
}
func (s *frpClientConfigProtocolImpl) Init(handler *action.Handler) error {
	if _, err := handler.SettingsInit(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to set default protocol: %w", err)
	}
	return nil
}
func (s *frpClientConfigProtocolImpl) Set(h *action.Handler, v any) error {
	str, err := config.GetConfigValString(v)
	if err != nil {
		return err
	}
	if str != "tcp" && str != "kcp" && str != "quic" && str != "websocket" {
		return errors.New("protocol must be one of TCP, KCP, QUIC or WebSocket")
	}
	return h.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), str)
}
func (s *frpClientConfigProtocolImpl) Get(h *action.Handler) (any, error) {
	val, _, _, err := h.SettingsGet(s.GetNamespace(), s.GetKey())
	if err != nil {
		return nil, fmt.Errorf("failed to get protocol: %w", err)
	}
	str, ok := val.(string)
	if !ok {
		return nil, errors.New("string expected")
	}
	return str, nil
}
func (s *frpClientConfigProtocolImpl) Restore(h *action.Handler) error {
	return h.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue())
}

type frpClientConfigProxyNameImpl struct{}

func (s *frpClientConfigProxyNameImpl) GetName() string             { return "Proxy Name" }
func (s *frpClientConfigProxyNameImpl) GetNamespace() string        { return ID }
func (s *frpClientConfigProxyNameImpl) GetKey() string              { return "proxy_name" }
func (s *frpClientConfigProxyNameImpl) GetType() action.SettingType { return action.String }
func (s *frpClientConfigProxyNameImpl) IsRequired() bool            { return true }
func (s *frpClientConfigProxyNameImpl) GetVersion() int             { return 0 }
func (s *frpClientConfigProxyNameImpl) GetOptions() map[string]any  { return nil }
func (s *frpClientConfigProxyNameImpl) GetDefaultValue() any {
	b := make([]byte, 4)
	if _, err := crypto_rand.Read(b); err != nil {
		return "anyshake-observer-web"
	}
	return fmt.Sprintf("anyshake-observer-%x", b)
}
func (s *frpClientConfigProxyNameImpl) GetDescription() string {
	return "Unique name for the proxy instance."
}
func (s *frpClientConfigProxyNameImpl) Init(handler *action.Handler) error {
	if _, err := handler.SettingsInit(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to set default proxy_name: %w", err)
	}
	return nil
}
func (s *frpClientConfigProxyNameImpl) Set(handler *action.Handler, newVal any) error {
	val, err := config.GetConfigValString(newVal)
	if err != nil {
		return err
	}
	if val == "" {
		return errors.New("proxy name cannot be empty")
	}
	return handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), val)
}
func (s *frpClientConfigProxyNameImpl) Get(handler *action.Handler) (any, error) {
	val, _, _, err := handler.SettingsGet(s.GetNamespace(), s.GetKey())
	if err != nil {
		return nil, fmt.Errorf("failed to get proxy_name: %w", err)
	}
	proxyName, ok := val.(string)
	if !ok {
		return nil, errors.New("string expected")
	}
	return proxyName, nil
}
func (s *frpClientConfigProxyNameImpl) Restore(handler *action.Handler) error {
	return handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue())
}

type frpClientConfigUseEncryptionImpl struct{}

func (s *frpClientConfigUseEncryptionImpl) GetName() string             { return "Use Encryption" }
func (s *frpClientConfigUseEncryptionImpl) GetNamespace() string        { return ID }
func (s *frpClientConfigUseEncryptionImpl) GetKey() string              { return "use_encryption" }
func (s *frpClientConfigUseEncryptionImpl) GetType() action.SettingType { return action.Bool }
func (s *frpClientConfigUseEncryptionImpl) IsRequired() bool            { return false }
func (s *frpClientConfigUseEncryptionImpl) GetVersion() int             { return 0 }
func (s *frpClientConfigUseEncryptionImpl) GetOptions() map[string]any  { return nil }
func (s *frpClientConfigUseEncryptionImpl) GetDefaultValue() any        { return true }
func (s *frpClientConfigUseEncryptionImpl) GetDescription() string {
	return "Enable encryption for the proxy traffic."
}
func (s *frpClientConfigUseEncryptionImpl) Init(handler *action.Handler) error {
	if _, err := handler.SettingsInit(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to set default use_encryption: %w", err)
	}
	return nil
}
func (s *frpClientConfigUseEncryptionImpl) Set(handler *action.Handler, newVal any) error {
	val, err := config.GetConfigValBool(newVal)
	if err != nil {
		return err
	}
	return handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), val)
}
func (s *frpClientConfigUseEncryptionImpl) Get(handler *action.Handler) (any, error) {
	val, _, _, err := handler.SettingsGet(s.GetNamespace(), s.GetKey())
	if err != nil {
		return nil, fmt.Errorf("failed to get use_encryption: %w", err)
	}
	encrypt, ok := val.(bool)
	if !ok {
		return nil, errors.New("bool expected")
	}
	return encrypt, nil
}
func (s *frpClientConfigUseEncryptionImpl) Restore(handler *action.Handler) error {
	return handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue())
}

type frpClientConfigUseCompressionImpl struct{}

func (s *frpClientConfigUseCompressionImpl) GetName() string             { return "Use Compression" }
func (s *frpClientConfigUseCompressionImpl) GetNamespace() string        { return ID }
func (s *frpClientConfigUseCompressionImpl) GetKey() string              { return "use_compression" }
func (s *frpClientConfigUseCompressionImpl) GetType() action.SettingType { return action.Bool }
func (s *frpClientConfigUseCompressionImpl) IsRequired() bool            { return false }
func (s *frpClientConfigUseCompressionImpl) GetVersion() int             { return 0 }
func (s *frpClientConfigUseCompressionImpl) GetOptions() map[string]any  { return nil }
func (s *frpClientConfigUseCompressionImpl) GetDefaultValue() any        { return true }
func (s *frpClientConfigUseCompressionImpl) GetDescription() string {
	return "Enable compression for the proxy traffic."
}
func (s *frpClientConfigUseCompressionImpl) Init(handler *action.Handler) error {
	if _, err := handler.SettingsInit(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to set default use_compression: %w", err)
	}
	return nil
}
func (s *frpClientConfigUseCompressionImpl) Set(handler *action.Handler, newVal any) error {
	val, err := config.GetConfigValBool(newVal)
	if err != nil {
		return err
	}
	return handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), val)
}
func (s *frpClientConfigUseCompressionImpl) Get(handler *action.Handler) (any, error) {
	val, _, _, err := handler.SettingsGet(s.GetNamespace(), s.GetKey())
	if err != nil {
		return nil, fmt.Errorf("failed to get use_compression: %w", err)
	}
	compress, ok := val.(bool)
	if !ok {
		return nil, errors.New("bool expected")
	}
	return compress, nil
}
func (s *frpClientConfigUseCompressionImpl) Restore(handler *action.Handler) error {
	return handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue())
}

type frpClientConfigUseDomainAccessImpl struct{}

func (s *frpClientConfigUseDomainAccessImpl) GetName() string             { return "Use Domain Access" }
func (s *frpClientConfigUseDomainAccessImpl) GetNamespace() string        { return ID }
func (s *frpClientConfigUseDomainAccessImpl) GetKey() string              { return "use_domain_access" }
func (s *frpClientConfigUseDomainAccessImpl) GetType() action.SettingType { return action.Bool }
func (s *frpClientConfigUseDomainAccessImpl) IsRequired() bool            { return true }
func (s *frpClientConfigUseDomainAccessImpl) GetVersion() int             { return 0 }
func (s *frpClientConfigUseDomainAccessImpl) GetOptions() map[string]any  { return nil }
func (s *frpClientConfigUseDomainAccessImpl) GetDefaultValue() any        { return false }
func (s *frpClientConfigUseDomainAccessImpl) GetDescription() string {
	return "Enable domain access for the proxy traffic. If enabled, the proxy will be accessible via subdomain or custom domains."
}
func (s *frpClientConfigUseDomainAccessImpl) Init(handler *action.Handler) error {
	if _, err := handler.SettingsInit(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to set default use_domain_access: %w", err)
	}
	return nil
}
func (s *frpClientConfigUseDomainAccessImpl) Set(handler *action.Handler, newVal any) error {
	val, err := config.GetConfigValBool(newVal)
	if err != nil {
		return err
	}
	return handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), val)
}
func (s *frpClientConfigUseDomainAccessImpl) Get(handler *action.Handler) (any, error) {
	val, _, _, err := handler.SettingsGet(s.GetNamespace(), s.GetKey())
	if err != nil {
		return nil, fmt.Errorf("failed to get use_domain_access: %w", err)
	}
	domainAccess, ok := val.(bool)
	if !ok {
		return nil, errors.New("bool expected")
	}
	return domainAccess, nil
}
func (s *frpClientConfigUseDomainAccessImpl) Restore(handler *action.Handler) error {
	return handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue())
}

type frpClientConfigRemoteOutboundPortImpl struct{}

func (s *frpClientConfigRemoteOutboundPortImpl) GetName() string             { return "Remote Outbound Port" }
func (s *frpClientConfigRemoteOutboundPortImpl) GetNamespace() string        { return ID }
func (s *frpClientConfigRemoteOutboundPortImpl) GetKey() string              { return "remote_outbound_port" }
func (s *frpClientConfigRemoteOutboundPortImpl) GetType() action.SettingType { return action.Int }
func (s *frpClientConfigRemoteOutboundPortImpl) IsRequired() bool            { return false }
func (s *frpClientConfigRemoteOutboundPortImpl) GetVersion() int             { return 0 }
func (s *frpClientConfigRemoteOutboundPortImpl) GetOptions() map[string]any  { return nil }
func (s *frpClientConfigRemoteOutboundPortImpl) GetDefaultValue() any {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	return rand.Intn(30000) + 30000
}
func (s *frpClientConfigRemoteOutboundPortImpl) GetDescription() string {
	return "Remote outbound port for the proxy traffic, this field is available only when domain mode is disabled."
}
func (s *frpClientConfigRemoteOutboundPortImpl) Init(handler *action.Handler) error {
	if _, err := handler.SettingsInit(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to set default remote_outbound_port: %w", err)
	}
	return nil
}
func (s *frpClientConfigRemoteOutboundPortImpl) Set(handler *action.Handler, newVal any) error {
	val, err := config.GetConfigValInt64(newVal)
	if err != nil {
		return err
	}
	if val < 1 || val > 65535 {
		return errors.New("remote outbound port must be between 1 and 65535")
	}
	return handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), val)
}
func (s *frpClientConfigRemoteOutboundPortImpl) Get(handler *action.Handler) (any, error) {
	val, _, _, err := handler.SettingsGet(s.GetNamespace(), s.GetKey())
	if err != nil {
		return nil, fmt.Errorf("failed to get remote_outbound_port: %w", err)
	}
	remoteOutboundPort, ok := val.(int64)
	if !ok {
		return nil, errors.New("int expected")
	}
	return remoteOutboundPort, nil
}
func (s *frpClientConfigRemoteOutboundPortImpl) Restore(handler *action.Handler) error {
	return handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue())
}

type frpClientConfigSubdomainImpl struct{}

func (s *frpClientConfigSubdomainImpl) GetName() string             { return "Subdomain" }
func (s *frpClientConfigSubdomainImpl) GetNamespace() string        { return ID }
func (s *frpClientConfigSubdomainImpl) GetKey() string              { return "subdomain" }
func (s *frpClientConfigSubdomainImpl) GetType() action.SettingType { return action.String }
func (s *frpClientConfigSubdomainImpl) IsRequired() bool            { return false }
func (s *frpClientConfigSubdomainImpl) GetVersion() int             { return 0 }
func (s *frpClientConfigSubdomainImpl) GetOptions() map[string]any  { return nil }
func (s *frpClientConfigSubdomainImpl) GetDefaultValue() any        { return "" }
func (s *frpClientConfigSubdomainImpl) GetDescription() string {
	return "Subdomain to bind for this proxy, this field is available only when domain mode is enabled."
}
func (s *frpClientConfigSubdomainImpl) Init(handler *action.Handler) error {
	if _, err := handler.SettingsInit(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to set default subdomain: %w", err)
	}
	return nil
}
func (s *frpClientConfigSubdomainImpl) Set(handler *action.Handler, newVal any) error {
	val, err := config.GetConfigValString(newVal)
	if err != nil {
		return err
	}
	return handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), val)
}
func (s *frpClientConfigSubdomainImpl) Get(handler *action.Handler) (any, error) {
	val, _, _, err := handler.SettingsGet(s.GetNamespace(), s.GetKey())
	if err != nil {
		return nil, fmt.Errorf("failed to get subdomain: %w", err)
	}
	subdomain, ok := val.(string)
	if !ok {
		return nil, errors.New("string expected")
	}
	return subdomain, nil
}
func (s *frpClientConfigSubdomainImpl) Restore(handler *action.Handler) error {
	return handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue())
}

type frpClientConfigCustomDomainsImpl struct{}

func (s *frpClientConfigCustomDomainsImpl) GetName() string             { return "Custom Domains" }
func (s *frpClientConfigCustomDomainsImpl) GetNamespace() string        { return ID }
func (s *frpClientConfigCustomDomainsImpl) GetKey() string              { return "custom_domains" }
func (s *frpClientConfigCustomDomainsImpl) GetType() action.SettingType { return action.StringArray }
func (s *frpClientConfigCustomDomainsImpl) IsRequired() bool            { return false }
func (s *frpClientConfigCustomDomainsImpl) GetVersion() int             { return 0 }
func (s *frpClientConfigCustomDomainsImpl) GetOptions() map[string]any  { return nil }
func (s *frpClientConfigCustomDomainsImpl) GetDefaultValue() any        { return []string{} }
func (s *frpClientConfigCustomDomainsImpl) GetDescription() string {
	return "List of custom domains to bind for this proxy, this field is available only when domain mode is enabled."
}
func (s *frpClientConfigCustomDomainsImpl) Init(handler *action.Handler) error {
	if _, err := handler.SettingsInit(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to set default custom_domains: %w", err)
	}
	return nil
}
func (s *frpClientConfigCustomDomainsImpl) Set(handler *action.Handler, newVal any) error {
	val, err := config.GetConfigValStringArray(newVal)
	if err != nil {
		return err
	}
	return handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), val)
}
func (s *frpClientConfigCustomDomainsImpl) Get(handler *action.Handler) (any, error) {
	val, _, _, err := handler.SettingsGet(s.GetNamespace(), s.GetKey())
	if err != nil {
		return nil, fmt.Errorf("failed to get custom_domains: %w", err)
	}
	customDomains, ok := val.([]string)
	if !ok {
		return nil, errors.New("string array expected")
	}
	return customDomains, nil
}
func (s *frpClientConfigCustomDomainsImpl) Restore(handler *action.Handler) error {
	return handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue())
}

func (s *FrpClientServiceImpl) GetConfigConstraint() []config.IConstraint {
	return []config.IConstraint{
		&frpClientConfigEnabledImpl{},
		&frpClientConfigServerAddrImpl{},
		&frpClientConfigServerPortImpl{},
		&frpClientConfigDisableCustomTLSFirstByteImpl{},
		&frpClientConfigUserImpl{},
		&frpClientConfigTokenImpl{},
		&frpClientConfigPoolCountImpl{},
		&frpClientConfigTcpMuxImpl{},
		&frpClientConfigTlsEnableImpl{},
		&frpClientConfigProtocolImpl{},
		&frpClientConfigProxyNameImpl{},
		&frpClientConfigUseEncryptionImpl{},
		&frpClientConfigUseCompressionImpl{},
		&frpClientConfigUseDomainAccessImpl{},
		&frpClientConfigRemoteOutboundPortImpl{},
		&frpClientConfigSubdomainImpl{},
		&frpClientConfigCustomDomainsImpl{},
	}
}
