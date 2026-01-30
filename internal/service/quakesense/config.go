package quakesense

import (
	"errors"
	"fmt"
	"time"

	crypto_rand "crypto/rand"

	"github.com/anyshake/observer/config"
	"github.com/anyshake/observer/internal/dao/action"
	"github.com/samber/lo"
)

type quakeSenseConfigEnabledImpl struct{}

func (s *quakeSenseConfigEnabledImpl) GetName() string             { return "Enable" }
func (s *quakeSenseConfigEnabledImpl) GetNamespace() string        { return ID }
func (s *quakeSenseConfigEnabledImpl) GetKey() string              { return "enabled" }
func (s *quakeSenseConfigEnabledImpl) GetType() action.SettingType { return action.Bool }
func (s *quakeSenseConfigEnabledImpl) IsRequired() bool            { return true }
func (s *quakeSenseConfigEnabledImpl) GetVersion() int             { return 0 }
func (s *quakeSenseConfigEnabledImpl) GetOptions() map[string]any  { return nil }
func (s *quakeSenseConfigEnabledImpl) GetDefaultValue() any        { return false }
func (s *quakeSenseConfigEnabledImpl) GetDescription() string {
	return "Enable QuakeSense service to detect realtime earthquakes, allows sending P-wave alerts via MQTT protocol"
}
func (s *quakeSenseConfigEnabledImpl) Init(handler *action.Handler) error {
	if _, err := handler.SettingsInit(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to set default QuakeSense service availability: %w", err)
	}
	return nil
}
func (s *quakeSenseConfigEnabledImpl) Set(handler *action.Handler, newVal any) error {
	enabled, err := config.GetConfigValBool(newVal)
	if err != nil {
		return err
	}
	if err := handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), enabled); err != nil {
		return fmt.Errorf("failed to set QuakeSense service availability: %w", err)
	}
	return nil
}
func (s *quakeSenseConfigEnabledImpl) Get(handler *action.Handler) (any, error) {
	val, _, _, err := handler.SettingsGet(s.GetNamespace(), s.GetKey())
	if err != nil {
		return nil, fmt.Errorf("failed to get QuakeSense service availability: %w", err)
	}
	enabled, ok := val.(bool)
	if !ok {
		return nil, errors.New("boolean expected")
	}
	return enabled, nil
}
func (s *quakeSenseConfigEnabledImpl) Restore(handler *action.Handler) error {
	if err := handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to reset QuakeSense service availability: %w", err)
	}
	return nil
}

type quakeSenseConfigMqttBrokerImpl struct{}

func (s *quakeSenseConfigMqttBrokerImpl) GetName() string             { return "MQTT Broker" }
func (s *quakeSenseConfigMqttBrokerImpl) GetNamespace() string        { return ID }
func (s *quakeSenseConfigMqttBrokerImpl) GetKey() string              { return "mqtt_broker" }
func (s *quakeSenseConfigMqttBrokerImpl) GetType() action.SettingType { return action.String }
func (s *quakeSenseConfigMqttBrokerImpl) IsRequired() bool            { return true }
func (s *quakeSenseConfigMqttBrokerImpl) GetVersion() int             { return 0 }
func (s *quakeSenseConfigMqttBrokerImpl) GetOptions() map[string]any  { return nil }
func (s *quakeSenseConfigMqttBrokerImpl) GetDefaultValue() any        { return "tcp://broker.emqx.io:1883" }
func (s *quakeSenseConfigMqttBrokerImpl) GetDescription() string {
	return "Address used for connecting to MQTT broker, by default, the host is tcp://broker.emqx.io:1883"
}
func (s *quakeSenseConfigMqttBrokerImpl) Init(handler *action.Handler) error {
	if _, err := handler.SettingsInit(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to set default MQTT broker: %w", err)
	}
	return nil
}
func (s *quakeSenseConfigMqttBrokerImpl) Set(handler *action.Handler, newVal any) error {
	broker, err := config.GetConfigValString(newVal)
	if err != nil {
		return err
	}
	if err := handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), broker); err != nil {
		return fmt.Errorf("failed to set MQTT broker: %w", err)
	}
	return nil
}
func (s *quakeSenseConfigMqttBrokerImpl) Get(handler *action.Handler) (any, error) {
	val, _, _, err := handler.SettingsGet(s.GetNamespace(), s.GetKey())
	if err != nil {
		return nil, fmt.Errorf("failed to get MQTT broker: %w", err)
	}
	broker, ok := val.(string)
	if !ok {
		return nil, errors.New("string expected")
	}
	return broker, nil
}
func (s *quakeSenseConfigMqttBrokerImpl) Restore(handler *action.Handler) error {
	if err := handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to reset MQTT broker: %w", err)
	}
	return nil
}

type quakeSenseConfigMqttTopicImpl struct{}

func (s *quakeSenseConfigMqttTopicImpl) GetName() string             { return "MQTT Topic" }
func (s *quakeSenseConfigMqttTopicImpl) GetNamespace() string        { return ID }
func (s *quakeSenseConfigMqttTopicImpl) GetKey() string              { return "mqtt_topic" }
func (s *quakeSenseConfigMqttTopicImpl) GetType() action.SettingType { return action.String }
func (s *quakeSenseConfigMqttTopicImpl) IsRequired() bool            { return true }
func (s *quakeSenseConfigMqttTopicImpl) GetVersion() int             { return 0 }
func (s *quakeSenseConfigMqttTopicImpl) GetOptions() map[string]any  { return nil }
func (s *quakeSenseConfigMqttTopicImpl) GetDefaultValue() any {
	b := make([]byte, 4)
	if _, err := crypto_rand.Read(b); err != nil {
		return fmt.Sprintf("anyshake/quakesense/%x", uint32(time.Now().UnixNano()))
	}
	return fmt.Sprintf("anyshake/quakesense/%x", b)
}
func (s *quakeSenseConfigMqttTopicImpl) GetDescription() string {
	return "Topic used for publishing MQTT messages, by default, the topic is anyshake/quakesense/<random string>"
}
func (s *quakeSenseConfigMqttTopicImpl) Init(handler *action.Handler) error {
	if _, err := handler.SettingsInit(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to set default MQTT topic: %w", err)
	}
	return nil
}
func (s *quakeSenseConfigMqttTopicImpl) Set(handler *action.Handler, newVal any) error {
	topic, err := config.GetConfigValString(newVal)
	if err != nil {
		return err
	}
	if err := handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), topic); err != nil {
		return fmt.Errorf("failed to set MQTT topic: %w", err)
	}
	return nil
}
func (s *quakeSenseConfigMqttTopicImpl) Get(handler *action.Handler) (any, error) {
	val, _, _, err := handler.SettingsGet(s.GetNamespace(), s.GetKey())
	if err != nil {
		return nil, fmt.Errorf("failed to get MQTT topic: %w", err)
	}
	topic, ok := val.(string)
	if !ok {
		return nil, errors.New("string expected")
	}
	return topic, nil
}
func (s *quakeSenseConfigMqttTopicImpl) Restore(handler *action.Handler) error {
	if err := handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to reset MQTT topic: %w", err)
	}
	return nil
}

type quakeSenseConfigMqttUsernameImpl struct{}

func (s *quakeSenseConfigMqttUsernameImpl) GetName() string             { return "MQTT Username" }
func (s *quakeSenseConfigMqttUsernameImpl) GetNamespace() string        { return ID }
func (s *quakeSenseConfigMqttUsernameImpl) GetKey() string              { return "mqtt_username" }
func (s *quakeSenseConfigMqttUsernameImpl) GetType() action.SettingType { return action.String }
func (s *quakeSenseConfigMqttUsernameImpl) IsRequired() bool            { return false }
func (s *quakeSenseConfigMqttUsernameImpl) GetVersion() int             { return 0 }
func (s *quakeSenseConfigMqttUsernameImpl) GetOptions() map[string]any  { return nil }
func (s *quakeSenseConfigMqttUsernameImpl) GetDefaultValue() any        { return "" }
func (s *quakeSenseConfigMqttUsernameImpl) GetDescription() string {
	return "MQTT Username used for connecting to MQTT broker, by default, the username is empty"
}
func (s *quakeSenseConfigMqttUsernameImpl) Init(handler *action.Handler) error {
	if _, err := handler.SettingsInit(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to set default MQTT username: %w", err)
	}
	return nil
}
func (s *quakeSenseConfigMqttUsernameImpl) Set(handler *action.Handler, newVal any) error {
	username, err := config.GetConfigValString(newVal)
	if err != nil {
		return err
	}
	if err := handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), username); err != nil {
		return fmt.Errorf("failed to set MQTT username: %w", err)
	}
	return nil
}
func (s *quakeSenseConfigMqttUsernameImpl) Get(handler *action.Handler) (any, error) {
	val, _, _, err := handler.SettingsGet(s.GetNamespace(), s.GetKey())
	if err != nil {
		return nil, fmt.Errorf("failed to get MQTT username: %w", err)
	}
	username, ok := val.(string)
	if !ok {
		return nil, errors.New("string expected")
	}
	return username, nil
}
func (s *quakeSenseConfigMqttUsernameImpl) Restore(handler *action.Handler) error {
	if err := handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to reset MQTT username: %w", err)
	}
	return nil
}

type quakeSenseConfigMqttPasswordImpl struct{}

func (s *quakeSenseConfigMqttPasswordImpl) GetName() string             { return "MQTT Password" }
func (s *quakeSenseConfigMqttPasswordImpl) GetNamespace() string        { return ID }
func (s *quakeSenseConfigMqttPasswordImpl) GetKey() string              { return "mqtt_password" }
func (s *quakeSenseConfigMqttPasswordImpl) GetType() action.SettingType { return action.String }
func (s *quakeSenseConfigMqttPasswordImpl) IsRequired() bool            { return false }
func (s *quakeSenseConfigMqttPasswordImpl) GetVersion() int             { return 0 }
func (s *quakeSenseConfigMqttPasswordImpl) GetOptions() map[string]any  { return nil }
func (s *quakeSenseConfigMqttPasswordImpl) GetDefaultValue() any        { return "" }
func (s *quakeSenseConfigMqttPasswordImpl) GetDescription() string {
	return "MQTT Password used for connecting to MQTT broker, by default, the password is empty"
}
func (s *quakeSenseConfigMqttPasswordImpl) Init(handler *action.Handler) error {
	if _, err := handler.SettingsInit(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to set default MQTT password: %w", err)
	}
	return nil
}
func (s *quakeSenseConfigMqttPasswordImpl) Set(handler *action.Handler, newVal any) error {
	password, err := config.GetConfigValString(newVal)
	if err != nil {
		return err
	}
	if err := handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), password); err != nil {
		return fmt.Errorf("failed to set MQTT password: %w", err)
	}
	return nil
}
func (s *quakeSenseConfigMqttPasswordImpl) Get(handler *action.Handler) (any, error) {
	val, _, _, err := handler.SettingsGet(s.GetNamespace(), s.GetKey())
	if err != nil {
		return nil, fmt.Errorf("failed to get MQTT password: %w", err)
	}
	password, ok := val.(string)
	if !ok {
		return nil, errors.New("string expected")
	}
	return password, nil
}
func (s *quakeSenseConfigMqttPasswordImpl) Restore(handler *action.Handler) error {
	if err := handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to reset MQTT password: %w", err)
	}
	return nil
}

type quakeSenseConfigMqttClientIdImpl struct{}

func (s *quakeSenseConfigMqttClientIdImpl) GetName() string             { return "MQTT Client ID" }
func (s *quakeSenseConfigMqttClientIdImpl) GetNamespace() string        { return ID }
func (s *quakeSenseConfigMqttClientIdImpl) GetKey() string              { return "mqtt_client_id" }
func (s *quakeSenseConfigMqttClientIdImpl) GetType() action.SettingType { return action.String }
func (s *quakeSenseConfigMqttClientIdImpl) IsRequired() bool            { return true }
func (s *quakeSenseConfigMqttClientIdImpl) GetVersion() int             { return 0 }
func (s *quakeSenseConfigMqttClientIdImpl) GetOptions() map[string]any  { return nil }
func (s *quakeSenseConfigMqttClientIdImpl) GetDefaultValue() any {
	b := make([]byte, 4)
	if _, err := crypto_rand.Read(b); err != nil {
		return fmt.Sprintf("anyshake-observer-%x", uint32(time.Now().UnixNano()))
	}
	return fmt.Sprintf("anyshake-observer-%x", b)
}
func (s *quakeSenseConfigMqttClientIdImpl) GetDescription() string {
	return "MQTT Client ID uniquely identifies a client when connecting to the broker. By default, it is anyshake-observer-<random string>."
}
func (s *quakeSenseConfigMqttClientIdImpl) Init(handler *action.Handler) error {
	if _, err := handler.SettingsInit(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to set default MQTT client ID: %w", err)
	}
	return nil
}
func (s *quakeSenseConfigMqttClientIdImpl) Set(handler *action.Handler, newVal any) error {
	topic, err := config.GetConfigValString(newVal)
	if err != nil {
		return err
	}
	if err := handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), topic); err != nil {
		return fmt.Errorf("failed to set MQTT client ID: %w", err)
	}
	return nil
}
func (s *quakeSenseConfigMqttClientIdImpl) Get(handler *action.Handler) (any, error) {
	val, _, _, err := handler.SettingsGet(s.GetNamespace(), s.GetKey())
	if err != nil {
		return nil, fmt.Errorf("failed to get MQTT client ID: %w", err)
	}
	id, ok := val.(string)
	if !ok {
		return nil, errors.New("string expected")
	}
	return id, nil
}
func (s *quakeSenseConfigMqttClientIdImpl) Restore(handler *action.Handler) error {
	if err := handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to reset MQTT client ID: %w", err)
	}
	return nil
}

type quakeSenseConfigMonitorChannelImpl struct{}

func (s *quakeSenseConfigMonitorChannelImpl) GetName() string             { return "Monitor Channel" }
func (s *quakeSenseConfigMonitorChannelImpl) GetNamespace() string        { return ID }
func (s *quakeSenseConfigMonitorChannelImpl) GetKey() string              { return "monitor_channel" }
func (s *quakeSenseConfigMonitorChannelImpl) GetType() action.SettingType { return action.String }
func (s *quakeSenseConfigMonitorChannelImpl) IsRequired() bool            { return true }
func (s *quakeSenseConfigMonitorChannelImpl) GetVersion() int             { return 0 }
func (s *quakeSenseConfigMonitorChannelImpl) GetOptions() map[string]any  { return nil }
func (s *quakeSenseConfigMonitorChannelImpl) GetDefaultValue() any        { return "EHZ" }
func (s *quakeSenseConfigMonitorChannelImpl) GetDescription() string {
	return "Specify which channel to monitor, by default, the channel is EHZ, it should be the same as the channel codes in station config"
}
func (s *quakeSenseConfigMonitorChannelImpl) Init(handler *action.Handler) error {
	if _, err := handler.SettingsInit(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to set default monitor channel: %w", err)
	}
	return nil
}
func (s *quakeSenseConfigMonitorChannelImpl) Set(handler *action.Handler, newVal any) error {
	channel, err := config.GetConfigValString(newVal)
	if err != nil {
		return err
	}
	if err := handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), channel); err != nil {
		return fmt.Errorf("failed to set monitor channel: %w", err)
	}
	return nil
}
func (s *quakeSenseConfigMonitorChannelImpl) Get(handler *action.Handler) (any, error) {
	val, _, _, err := handler.SettingsGet(s.GetNamespace(), s.GetKey())
	if err != nil {
		return nil, fmt.Errorf("failed to get monitor channel: %w", err)
	}
	channel, ok := val.(string)
	if !ok {
		return nil, errors.New("string expected")
	}
	return channel, nil
}
func (s *quakeSenseConfigMonitorChannelImpl) Restore(handler *action.Handler) error {
	if err := handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to reset monitor channel: %w", err)
	}
	return nil
}

type quakeSenseConfigFilterTypeImpl struct{}

func (s *quakeSenseConfigFilterTypeImpl) GetName() string             { return "Filter Type" }
func (s *quakeSenseConfigFilterTypeImpl) GetNamespace() string        { return ID }
func (s *quakeSenseConfigFilterTypeImpl) GetKey() string              { return "filter_type" }
func (s *quakeSenseConfigFilterTypeImpl) GetType() action.SettingType { return action.String }
func (s *quakeSenseConfigFilterTypeImpl) IsRequired() bool            { return true }
func (s *quakeSenseConfigFilterTypeImpl) GetVersion() int             { return 0 }
func (s *quakeSenseConfigFilterTypeImpl) GetOptions() map[string]any {
	return map[string]any{
		"No Filter":              NO_FILTER,
		"Low-pass Filter (LPF)":  LOW_PASS_FILTER,
		"Band-pass Filter (BPF)": BAND_PASS_FILTER,
		"High-pass Filter (HPF)": HIGH_PASS_FILTER,
	}
}
func (s *quakeSenseConfigFilterTypeImpl) GetDefaultValue() any { return BAND_PASS_FILTER }
func (s *quakeSenseConfigFilterTypeImpl) GetDescription() string {
	return "Specify which filter type to use, by default, the filter type is Band-pass filter."
}
func (s *quakeSenseConfigFilterTypeImpl) Init(handler *action.Handler) error {
	if _, err := handler.SettingsInit(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to set default filter type: %w", err)
	}
	return nil
}
func (s *quakeSenseConfigFilterTypeImpl) Set(handler *action.Handler, newVal any) error {
	filterType, err := config.GetConfigValString(newVal)
	if err != nil {
		return err
	}
	if !lo.Contains([]string{
		NO_FILTER,
		LOW_PASS_FILTER,
		BAND_PASS_FILTER,
		HIGH_PASS_FILTER,
	}, filterType) {
		return fmt.Errorf("invalid filter type: %s", filterType)
	}
	if err := handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), filterType); err != nil {
		return fmt.Errorf("failed to set filter type: %w", err)
	}
	return nil
}
func (s *quakeSenseConfigFilterTypeImpl) Get(handler *action.Handler) (any, error) {
	val, _, _, err := handler.SettingsGet(s.GetNamespace(), s.GetKey())
	if err != nil {
		return nil, fmt.Errorf("failed to get filter type: %w", err)
	}
	filterType, ok := val.(string)
	if !ok {
		return nil, errors.New("string expected")
	}
	return filterType, nil
}
func (s *quakeSenseConfigFilterTypeImpl) Restore(handler *action.Handler) error {
	if err := handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to reset filter type: %w", err)
	}
	return nil
}

type quakeSenseConfigMinFreqImpl struct{}

func (s *quakeSenseConfigMinFreqImpl) GetName() string             { return "Minimum Frequency" }
func (s *quakeSenseConfigMinFreqImpl) GetNamespace() string        { return ID }
func (s *quakeSenseConfigMinFreqImpl) GetKey() string              { return "min_freq" }
func (s *quakeSenseConfigMinFreqImpl) GetType() action.SettingType { return action.Float }
func (s *quakeSenseConfigMinFreqImpl) IsRequired() bool            { return true }
func (s *quakeSenseConfigMinFreqImpl) GetOptions() map[string]any  { return nil }
func (s *quakeSenseConfigMinFreqImpl) GetVersion() int             { return 0 }
func (s *quakeSenseConfigMinFreqImpl) GetDefaultValue() any        { return 0.5 }
func (s *quakeSenseConfigMinFreqImpl) GetDescription() string {
	return "Specify the minimum frequency to filter, by default, the minimum frequency is 0.5 Hz."
}
func (s *quakeSenseConfigMinFreqImpl) Init(handler *action.Handler) error {
	if _, err := handler.SettingsInit(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to set default minimum frequency: %w", err)
	}
	return nil
}
func (s *quakeSenseConfigMinFreqImpl) Set(handler *action.Handler, newVal any) error {
	minFreq, err := config.GetConfigValFloat64(newVal)
	if err != nil {
		return err
	}
	if minFreq < 0 {
		return errors.New("minimum frequency cannot be negative")
	}
	if err := handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), minFreq); err != nil {
		return fmt.Errorf("failed to set minimum frequency: %w", err)
	}
	return nil
}
func (s *quakeSenseConfigMinFreqImpl) Get(handler *action.Handler) (any, error) {
	val, _, _, err := handler.SettingsGet(s.GetNamespace(), s.GetKey())
	if err != nil {
		return nil, fmt.Errorf("failed to get minimum frequency: %w", err)
	}
	minFreq, ok := val.(float64)
	if !ok {
		return nil, errors.New("float64 expected")
	}
	return minFreq, nil
}
func (s *quakeSenseConfigMinFreqImpl) Restore(handler *action.Handler) error {
	if err := handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to reset minimum frequency: %w", err)
	}
	return nil
}

type quakeSenseConfigMaxFreqImpl struct{}

func (s *quakeSenseConfigMaxFreqImpl) GetName() string             { return "Maximum Frequency" }
func (s *quakeSenseConfigMaxFreqImpl) GetNamespace() string        { return ID }
func (s *quakeSenseConfigMaxFreqImpl) GetKey() string              { return "max_freq" }
func (s *quakeSenseConfigMaxFreqImpl) GetType() action.SettingType { return action.Float }
func (s *quakeSenseConfigMaxFreqImpl) IsRequired() bool            { return true }
func (s *quakeSenseConfigMaxFreqImpl) GetOptions() map[string]any  { return nil }
func (s *quakeSenseConfigMaxFreqImpl) GetVersion() int             { return 0 }
func (s *quakeSenseConfigMaxFreqImpl) GetDefaultValue() any        { return 10.0 }
func (s *quakeSenseConfigMaxFreqImpl) GetDescription() string {
	return "Specify the maximum frequency to filter, by default, the maximum frequency is 10 Hz."
}
func (s *quakeSenseConfigMaxFreqImpl) Init(handler *action.Handler) error {
	if _, err := handler.SettingsInit(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to set default maximum frequency: %w", err)
	}
	return nil
}
func (s *quakeSenseConfigMaxFreqImpl) Set(handler *action.Handler, newVal any) error {
	maxFreq, err := config.GetConfigValFloat64(newVal)
	if err != nil {
		return err
	}
	if maxFreq < 0 {
		return errors.New("maximum frequency cannot be negative")
	}
	if err := handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), maxFreq); err != nil {
		return fmt.Errorf("failed to set maximum frequency: %w", err)
	}
	return nil
}
func (s *quakeSenseConfigMaxFreqImpl) Get(handler *action.Handler) (any, error) {
	val, _, _, err := handler.SettingsGet(s.GetNamespace(), s.GetKey())
	if err != nil {
		return nil, fmt.Errorf("failed to get maximum frequency: %w", err)
	}
	maxFreq, ok := val.(float64)
	if !ok {
		return nil, errors.New("float64 expected")
	}
	return maxFreq, nil
}
func (s *quakeSenseConfigMaxFreqImpl) Restore(handler *action.Handler) error {
	if err := handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to reset maximum frequency: %w", err)
	}
	return nil
}

type quakeSenseConfigTriggerMethodImpl struct{}

func (s *quakeSenseConfigTriggerMethodImpl) GetName() string             { return "Trigger Method" }
func (s *quakeSenseConfigTriggerMethodImpl) GetNamespace() string        { return ID }
func (s *quakeSenseConfigTriggerMethodImpl) GetKey() string              { return "trigger_method" }
func (s *quakeSenseConfigTriggerMethodImpl) GetType() action.SettingType { return action.String }
func (s *quakeSenseConfigTriggerMethodImpl) IsRequired() bool            { return true }
func (s *quakeSenseConfigTriggerMethodImpl) GetOptions() map[string]any {
	return map[string]any{
		"Classic STA/LTA": CLASSIC_STA_LTA,
		// "Recursive STA/LTA": RECURSIVE_STA_LTA,
		// "Delayed STA/LTA":   DELAYED_STA_LTA,
		"Z-Detect": Z_DETECT,
	}
}
func (s *quakeSenseConfigTriggerMethodImpl) GetVersion() int      { return 0 }
func (s *quakeSenseConfigTriggerMethodImpl) GetDefaultValue() any { return CLASSIC_STA_LTA }
func (s *quakeSenseConfigTriggerMethodImpl) GetDescription() string {
	return "Specify the method to use for triggering, by default, the trigger method is Classic STA/LTA	."
}
func (s *quakeSenseConfigTriggerMethodImpl) Init(handler *action.Handler) error {
	if _, err := handler.SettingsInit(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to set default trigger method: %w", err)
	}
	return nil
}
func (s *quakeSenseConfigTriggerMethodImpl) Set(handler *action.Handler, newVal any) error {
	triggerMethod, err := config.GetConfigValString(newVal)
	if err != nil {
		return err
	}
	if !lo.Contains([]string{
		CLASSIC_STA_LTA,
		RECURSIVE_STA_LTA,
		DELAYED_STA_LTA,
		Z_DETECT,
	}, triggerMethod) {
		return fmt.Errorf("invalid trigger method: %s", triggerMethod)
	}
	if err := handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), triggerMethod); err != nil {
		return fmt.Errorf("failed to set trigger method: %w", err)
	}
	return nil
}
func (s *quakeSenseConfigTriggerMethodImpl) Get(handler *action.Handler) (any, error) {
	val, _, _, err := handler.SettingsGet(s.GetNamespace(), s.GetKey())
	if err != nil {
		return nil, fmt.Errorf("failed to get trigger method: %w", err)
	}
	method, ok := val.(string)
	if !ok {
		return nil, errors.New("string expected")
	}
	return method, nil
}
func (s *quakeSenseConfigTriggerMethodImpl) Restore(handler *action.Handler) error {
	if err := handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to reset trigger method: %w", err)
	}
	return nil
}

type quakeSenseConfigStaWindowImpl struct{}

func (s *quakeSenseConfigStaWindowImpl) GetName() string             { return "STA Window" }
func (s *quakeSenseConfigStaWindowImpl) GetNamespace() string        { return ID }
func (s *quakeSenseConfigStaWindowImpl) GetKey() string              { return "sta_window" }
func (s *quakeSenseConfigStaWindowImpl) GetType() action.SettingType { return action.Float }
func (s *quakeSenseConfigStaWindowImpl) IsRequired() bool            { return true }
func (s *quakeSenseConfigStaWindowImpl) GetOptions() map[string]any  { return nil }
func (s *quakeSenseConfigStaWindowImpl) GetVersion() int             { return 0 }
func (s *quakeSenseConfigStaWindowImpl) GetDefaultValue() any        { return 4.0 }
func (s *quakeSenseConfigStaWindowImpl) GetDescription() string {
	return "Specify the STA window, by default, the STA window is 4."
}
func (s *quakeSenseConfigStaWindowImpl) Init(handler *action.Handler) error {
	if _, err := handler.SettingsInit(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to set default STA window: %w", err)
	}
	return nil
}
func (s *quakeSenseConfigStaWindowImpl) Set(handler *action.Handler, newVal any) error {
	staWindow, err := config.GetConfigValFloat64(newVal)
	if err != nil {
		return err
	}
	if staWindow <= 0 {
		return errors.New("STA window cannot be zero or negative")
	}
	if staWindow > 300 {
		return errors.New("STA window cannot be greater than 300")
	}
	if err := handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), staWindow); err != nil {
		return fmt.Errorf("failed to set STA window: %w", err)
	}
	return nil
}
func (s *quakeSenseConfigStaWindowImpl) Get(handler *action.Handler) (any, error) {
	val, _, _, err := handler.SettingsGet(s.GetNamespace(), s.GetKey())
	if err != nil {
		return nil, fmt.Errorf("failed to get STA window: %w", err)
	}
	staWindow, ok := val.(float64)
	if !ok {
		return nil, errors.New("float64 expected")
	}
	return staWindow, nil
}
func (s *quakeSenseConfigStaWindowImpl) Restore(handler *action.Handler) error {
	if err := handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to reset STA window: %w", err)
	}
	return nil
}

type quakeSenseConfigLtaWindowImpl struct{}

func (s *quakeSenseConfigLtaWindowImpl) GetName() string             { return "LTA Window" }
func (s *quakeSenseConfigLtaWindowImpl) GetNamespace() string        { return ID }
func (s *quakeSenseConfigLtaWindowImpl) GetKey() string              { return "lta_window" }
func (s *quakeSenseConfigLtaWindowImpl) GetType() action.SettingType { return action.Float }
func (s *quakeSenseConfigLtaWindowImpl) IsRequired() bool            { return true }
func (s *quakeSenseConfigLtaWindowImpl) GetOptions() map[string]any  { return nil }
func (s *quakeSenseConfigLtaWindowImpl) GetVersion() int             { return 0 }
func (s *quakeSenseConfigLtaWindowImpl) GetDefaultValue() any        { return 30.0 }
func (s *quakeSenseConfigLtaWindowImpl) GetDescription() string {
	return "Specify the LTA window, by default, the LTA window is 30."
}
func (s *quakeSenseConfigLtaWindowImpl) Init(handler *action.Handler) error {
	if _, err := handler.SettingsInit(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to set default LTA window: %w", err)
	}
	return nil
}
func (s *quakeSenseConfigLtaWindowImpl) Set(handler *action.Handler, newVal any) error {
	ltaWindow, err := config.GetConfigValFloat64(newVal)
	if err != nil {
		return err
	}
	if ltaWindow <= 0 {
		return errors.New("LTA window cannot be zero or negative")
	}
	if ltaWindow > 300 {
		return errors.New("LTA window cannot be greater than 300")
	}
	if err := handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), ltaWindow); err != nil {
		return fmt.Errorf("failed to set LTA window: %w", err)
	}
	return nil
}
func (s *quakeSenseConfigLtaWindowImpl) Get(handler *action.Handler) (any, error) {
	val, _, _, err := handler.SettingsGet(s.GetNamespace(), s.GetKey())
	if err != nil {
		return nil, fmt.Errorf("failed to get LTA window: %w", err)
	}
	ltaWindow, ok := val.(float64)
	if !ok {
		return nil, errors.New("float64 expected")
	}
	return ltaWindow, nil
}
func (s *quakeSenseConfigLtaWindowImpl) Restore(handler *action.Handler) error {
	if err := handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to reset LTA window: %w", err)
	}
	return nil
}

type quakeSenseConfigTrigOnImpl struct{}

func (s *quakeSenseConfigTrigOnImpl) GetName() string             { return "Trigger On" }
func (s *quakeSenseConfigTrigOnImpl) GetNamespace() string        { return ID }
func (s *quakeSenseConfigTrigOnImpl) GetKey() string              { return "trig_on" }
func (s *quakeSenseConfigTrigOnImpl) GetType() action.SettingType { return action.Float }
func (s *quakeSenseConfigTrigOnImpl) IsRequired() bool            { return true }
func (s *quakeSenseConfigTrigOnImpl) GetOptions() map[string]any  { return nil }
func (s *quakeSenseConfigTrigOnImpl) GetVersion() int             { return 0 }
func (s *quakeSenseConfigTrigOnImpl) GetDefaultValue() any        { return 1.0 }
func (s *quakeSenseConfigTrigOnImpl) GetDescription() string {
	return "Specify the STA/LTA ratio to turn the trigger on, by default, the trigger on value is 1."
}
func (s *quakeSenseConfigTrigOnImpl) Init(handler *action.Handler) error {
	if _, err := handler.SettingsInit(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to set default trigger on value: %w", err)
	}
	return nil
}
func (s *quakeSenseConfigTrigOnImpl) Set(handler *action.Handler, newVal any) error {
	trigOn, err := config.GetConfigValFloat64(newVal)
	if err != nil {
		return err
	}
	if trigOn <= 0 {
		return errors.New("trigger on value cannot be zero or negative")
	}
	if err := handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), trigOn); err != nil {
		return fmt.Errorf("failed to set trigger on value: %w", err)
	}
	return nil
}
func (s *quakeSenseConfigTrigOnImpl) Get(handler *action.Handler) (any, error) {
	val, _, _, err := handler.SettingsGet(s.GetNamespace(), s.GetKey())
	if err != nil {
		return nil, fmt.Errorf("failed to get trigger on value: %w", err)
	}
	trigOn, ok := val.(float64)
	if !ok {
		return nil, errors.New("float64 expected")
	}
	return trigOn, nil
}
func (s *quakeSenseConfigTrigOnImpl) Restore(handler *action.Handler) error {
	if err := handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to reset trigger on value: %w", err)
	}
	return nil
}

type quakeSenseConfigTrigOffImpl struct{}

func (s *quakeSenseConfigTrigOffImpl) GetName() string             { return "Trigger Off" }
func (s *quakeSenseConfigTrigOffImpl) GetNamespace() string        { return ID }
func (s *quakeSenseConfigTrigOffImpl) GetKey() string              { return "trig_off" }
func (s *quakeSenseConfigTrigOffImpl) GetType() action.SettingType { return action.Float }
func (s *quakeSenseConfigTrigOffImpl) IsRequired() bool            { return true }
func (s *quakeSenseConfigTrigOffImpl) GetOptions() map[string]any  { return nil }
func (s *quakeSenseConfigTrigOffImpl) GetVersion() int             { return 0 }
func (s *quakeSenseConfigTrigOffImpl) GetDefaultValue() any        { return 0.4 }
func (s *quakeSenseConfigTrigOffImpl) GetDescription() string {
	return "Specify the STA/LTA ratio to turn the trigger off, by default, the trigger off value is 0.4."
}
func (s *quakeSenseConfigTrigOffImpl) Init(handler *action.Handler) error {
	if _, err := handler.SettingsInit(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to set default trigger off value: %w", err)
	}
	return nil
}
func (s *quakeSenseConfigTrigOffImpl) Set(handler *action.Handler, newVal any) error {
	trigOff, err := config.GetConfigValFloat64(newVal)
	if err != nil {
		return err
	}
	if trigOff <= 0 {
		return errors.New("trigger off value cannot be zero or negative")
	}
	if err := handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), trigOff); err != nil {
		return fmt.Errorf("failed to set trigger off value: %w", err)
	}
	return nil
}
func (s *quakeSenseConfigTrigOffImpl) Get(handler *action.Handler) (any, error) {
	val, _, _, err := handler.SettingsGet(s.GetNamespace(), s.GetKey())
	if err != nil {
		return nil, fmt.Errorf("failed to get trigger off value: %w", err)
	}
	trigOff, ok := val.(float64)
	if !ok {
		return nil, errors.New("float64 expected")
	}
	return trigOff, nil
}
func (s *quakeSenseConfigTrigOffImpl) Restore(handler *action.Handler) error {
	if err := handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to reset trigger off value: %w", err)
	}
	return nil
}

type quakeSenseConfigThrottleImpl struct{}

func (s *quakeSenseConfigThrottleImpl) GetName() string             { return "Throttle" }
func (s *quakeSenseConfigThrottleImpl) GetNamespace() string        { return ID }
func (s *quakeSenseConfigThrottleImpl) GetKey() string              { return "throttle_interval" }
func (s *quakeSenseConfigThrottleImpl) GetType() action.SettingType { return action.Int }
func (s *quakeSenseConfigThrottleImpl) IsRequired() bool            { return true }
func (s *quakeSenseConfigThrottleImpl) GetOptions() map[string]any  { return nil }
func (s *quakeSenseConfigThrottleImpl) GetVersion() int             { return 0 }
func (s *quakeSenseConfigThrottleImpl) GetDefaultValue() any        { return 10 }
func (s *quakeSenseConfigThrottleImpl) GetDescription() string {
	return "Specify the throttle interval in seconds, to disable throttling, set this field to 0, by default, the throttle interval is 10 seconds."
}
func (s *quakeSenseConfigThrottleImpl) Init(handler *action.Handler) error {
	if _, err := handler.SettingsInit(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to set default throttle interval value: %w", err)
	}
	return nil
}
func (s *quakeSenseConfigThrottleImpl) Set(handler *action.Handler, newVal any) error {
	throttle, err := config.GetConfigValInt64(newVal)
	if err != nil {
		return err
	}
	if throttle < 0 {
		return errors.New("throttle interval cannot be negative")
	}
	if err := handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), throttle); err != nil {
		return fmt.Errorf("failed to set throttle interval value: %w", err)
	}
	return nil
}
func (s *quakeSenseConfigThrottleImpl) Get(handler *action.Handler) (any, error) {
	val, _, _, err := handler.SettingsGet(s.GetNamespace(), s.GetKey())
	if err != nil {
		return nil, fmt.Errorf("failed to get throttle interval value: %w", err)
	}
	throttle, ok := val.(int64)
	if !ok {
		return nil, errors.New("int expected")
	}
	return throttle, nil
}
func (s *quakeSenseConfigThrottleImpl) Restore(handler *action.Handler) error {
	if err := handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to reset throttle interval value: %w", err)
	}
	return nil
}

func (s *QuakeSenseServiceImpl) GetConfigConstraint() []config.IConstraint {
	return []config.IConstraint{
		&quakeSenseConfigEnabledImpl{},
		&quakeSenseConfigMqttBrokerImpl{},
		&quakeSenseConfigMqttTopicImpl{},
		&quakeSenseConfigMqttUsernameImpl{},
		&quakeSenseConfigMqttPasswordImpl{},
		&quakeSenseConfigMqttClientIdImpl{},
		&quakeSenseConfigMonitorChannelImpl{},
		&quakeSenseConfigFilterTypeImpl{},
		&quakeSenseConfigMinFreqImpl{},
		&quakeSenseConfigMaxFreqImpl{},
		&quakeSenseConfigTriggerMethodImpl{},
		&quakeSenseConfigStaWindowImpl{},
		&quakeSenseConfigLtaWindowImpl{},
		&quakeSenseConfigTrigOnImpl{},
		&quakeSenseConfigTrigOffImpl{},
		&quakeSenseConfigThrottleImpl{},
	}
}
