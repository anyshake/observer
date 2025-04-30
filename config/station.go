package config

import (
	"errors"
	"fmt"
	"strings"

	"github.com/anyshake/observer/internal/dao/action"
)

const STATION_NAMESPACE = "global_station"

type StationNameConfigConstraintImpl struct{}

func (s *StationNameConfigConstraintImpl) GetName() string             { return "Station Name" }
func (s *StationNameConfigConstraintImpl) GetNamespace() string        { return STATION_NAMESPACE }
func (s *StationNameConfigConstraintImpl) GetKey() string              { return "station_name" }
func (s *StationNameConfigConstraintImpl) GetType() action.SettingType { return action.String }
func (s *StationNameConfigConstraintImpl) IsRequired() bool            { return true }
func (s *StationNameConfigConstraintImpl) GetVersion() int             { return 0 }
func (s *StationNameConfigConstraintImpl) GetOptions() map[string]any  { return nil }
func (s *StationNameConfigConstraintImpl) GetDefaultValue() any        { return "AnyShake Station" }
func (s *StationNameConfigConstraintImpl) GetDescription() string {
	return "Custom name for the station, should less than 20 characters."
}
func (s *StationNameConfigConstraintImpl) Init(handler *action.Handler) error {
	if _, err := handler.SettingsInit(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to set default station name: %w", err)
	}
	return nil
}
func (s *StationNameConfigConstraintImpl) Set(handler *action.Handler, newVal any) error {
	str, err := GetConfigValString(newVal)
	if err != nil {
		return err
	}
	if len(str) > 20 {
		return errors.New("station name must be less than 20 characters")
	}
	if err := handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), newVal); err != nil {
		return fmt.Errorf("failed to set station name: %w", err)
	}
	return nil
}
func (s *StationNameConfigConstraintImpl) Get(handler *action.Handler) (any, error) {
	val, _, _, err := handler.SettingsGet(s.GetNamespace(), s.GetKey())
	if err != nil {
		return nil, fmt.Errorf("failed to get station name: %w", err)
	}
	if _, ok := val.(string); !ok {
		return nil, errors.New("string expected")
	}
	return val, nil
}
func (s *StationNameConfigConstraintImpl) Restore(handler *action.Handler) error {
	if err := handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to reset station name: %w", err)
	}
	return nil
}

type StationDescriptionConfigConstraintImpl struct{}

func (s *StationDescriptionConfigConstraintImpl) GetName() string             { return "Station Description" }
func (s *StationDescriptionConfigConstraintImpl) GetNamespace() string        { return STATION_NAMESPACE }
func (s *StationDescriptionConfigConstraintImpl) GetKey() string              { return "station_description" }
func (s *StationDescriptionConfigConstraintImpl) GetType() action.SettingType { return action.String }
func (s *StationDescriptionConfigConstraintImpl) IsRequired() bool            { return false }
func (s *StationDescriptionConfigConstraintImpl) GetOptions() map[string]any  { return nil }
func (s *StationDescriptionConfigConstraintImpl) GetVersion() int             { return 0 }
func (s *StationDescriptionConfigConstraintImpl) GetDefaultValue() any {
	return "Proudly powered by AnyShake Project."
}
func (s *StationDescriptionConfigConstraintImpl) GetDescription() string {
	return "Brief description for the station, should less than 100 characters."
}
func (s *StationDescriptionConfigConstraintImpl) Init(handler *action.Handler) error {
	if _, err := handler.SettingsInit(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to set default station description: %w", err)
	}
	return nil
}
func (s *StationDescriptionConfigConstraintImpl) Set(handler *action.Handler, newVal any) error {
	str, err := GetConfigValString(newVal)
	if err != nil {
		return err
	}
	if len(str) > 100 {
		return errors.New("station description must be less than 20 characters")
	}
	if err := handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), newVal); err != nil {
		return fmt.Errorf("failed to set station description: %w", err)
	}
	return nil
}
func (s *StationDescriptionConfigConstraintImpl) Get(handler *action.Handler) (any, error) {
	val, _, _, err := handler.SettingsGet(s.GetNamespace(), s.GetKey())
	if err != nil {
		return nil, fmt.Errorf("failed to get station description: %w", err)
	}
	if _, ok := val.(string); !ok {
		return nil, errors.New("string expected")
	}
	return val, nil
}
func (s *StationDescriptionConfigConstraintImpl) Restore(handler *action.Handler) error {
	if err := handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to reset station description: %w", err)
	}
	return nil
}

type StationCountryConfigConstraintImpl struct{}

func (s *StationCountryConfigConstraintImpl) GetName() string             { return "Station Country" }
func (s *StationCountryConfigConstraintImpl) GetNamespace() string        { return STATION_NAMESPACE }
func (s *StationCountryConfigConstraintImpl) GetKey() string              { return "station_country" }
func (s *StationCountryConfigConstraintImpl) GetType() action.SettingType { return action.String }
func (s *StationCountryConfigConstraintImpl) IsRequired() bool            { return false }
func (s *StationCountryConfigConstraintImpl) GetVersion() int             { return 0 }
func (s *StationCountryConfigConstraintImpl) GetOptions() map[string]any  { return nil }
func (s *StationCountryConfigConstraintImpl) GetDefaultValue() any        { return "United States" }
func (s *StationCountryConfigConstraintImpl) GetDescription() string {
	return "Country where the station is located, should less than 20 characters, e.g. United States."
}
func (s *StationCountryConfigConstraintImpl) Init(handler *action.Handler) error {
	if _, err := handler.SettingsInit(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to set default station country: %w", err)
	}
	return nil
}
func (s *StationCountryConfigConstraintImpl) Set(handler *action.Handler, newVal any) error {
	str, err := GetConfigValString(newVal)
	if err != nil {
		return err
	}
	if len(str) > 20 {
		return errors.New("station country should be less than 20 characters")
	}
	if err := handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), strings.ToUpper(str)); err != nil {
		return fmt.Errorf("failed to set station country: %w", err)
	}
	return nil
}
func (s *StationCountryConfigConstraintImpl) Get(handler *action.Handler) (any, error) {
	val, _, _, err := handler.SettingsGet(s.GetNamespace(), s.GetKey())
	if err != nil {
		return nil, fmt.Errorf("failed to get station country: %w", err)
	}
	if _, ok := val.(string); !ok {
		return nil, errors.New("string expected")
	}
	return strings.ToUpper(val.(string)), nil
}
func (s *StationCountryConfigConstraintImpl) Restore(handler *action.Handler) error {
	if err := handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to reset station country: %w", err)
	}
	return nil
}

type StationPlaceConfigConstraintImpl struct{}

func (s *StationPlaceConfigConstraintImpl) GetName() string             { return "Station Place" }
func (s *StationPlaceConfigConstraintImpl) GetNamespace() string        { return STATION_NAMESPACE }
func (s *StationPlaceConfigConstraintImpl) GetKey() string              { return "station_place" }
func (s *StationPlaceConfigConstraintImpl) GetType() action.SettingType { return action.String }
func (s *StationPlaceConfigConstraintImpl) IsRequired() bool            { return false }
func (s *StationPlaceConfigConstraintImpl) GetVersion() int             { return 0 }
func (s *StationPlaceConfigConstraintImpl) GetOptions() map[string]any  { return nil }
func (s *StationPlaceConfigConstraintImpl) GetDefaultValue() any        { return "New York, NY" }
func (s *StationPlaceConfigConstraintImpl) GetDescription() string {
	return "Describe the location of the station, should include the city and less than 20 characters, e.g. New York, NY."
}
func (s *StationPlaceConfigConstraintImpl) Init(handler *action.Handler) error {
	if _, err := handler.SettingsInit(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to set default station place: %w", err)
	}
	return nil
}
func (s *StationPlaceConfigConstraintImpl) Set(handler *action.Handler, newVal any) error {
	str, err := GetConfigValString(newVal)
	if err != nil {
		return err
	}
	if len(str) > 20 {
		return errors.New("station place should be less than 20 characters")
	}
	if err := handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), str); err != nil {
		return fmt.Errorf("failed to set station place: %w", err)
	}
	return nil
}
func (s *StationPlaceConfigConstraintImpl) Get(handler *action.Handler) (any, error) {
	val, _, _, err := handler.SettingsGet(s.GetNamespace(), s.GetKey())
	if err != nil {
		return nil, fmt.Errorf("failed to get station place: %w", err)
	}

	if _, ok := val.(string); !ok {
		return nil, errors.New("string expected")
	}
	return val, nil
}

func (s *StationPlaceConfigConstraintImpl) Restore(handler *action.Handler) error {
	if err := handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to reset station place: %w", err)
	}
	return nil
}

type StationAffiliationConfigConstraintImpl struct{}

func (s *StationAffiliationConfigConstraintImpl) GetName() string             { return "Station Affiliation" }
func (s *StationAffiliationConfigConstraintImpl) GetNamespace() string        { return STATION_NAMESPACE }
func (s *StationAffiliationConfigConstraintImpl) GetKey() string              { return "station_affiliation" }
func (s *StationAffiliationConfigConstraintImpl) GetType() action.SettingType { return action.String }
func (s *StationAffiliationConfigConstraintImpl) IsRequired() bool            { return true }
func (s *StationAffiliationConfigConstraintImpl) GetVersion() int             { return 0 }
func (s *StationAffiliationConfigConstraintImpl) GetOptions() map[string]any  { return nil }
func (s *StationAffiliationConfigConstraintImpl) GetDefaultValue() any        { return "AnyShake Project" }
func (s *StationAffiliationConfigConstraintImpl) GetDescription() string {
	return "Who is the station affiliated with, should less than 50 characters, e.g. MIT Earth Resources Lab."
}
func (s *StationAffiliationConfigConstraintImpl) Init(handler *action.Handler) error {
	if _, err := handler.SettingsInit(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to set default station affiliation: %w", err)
	}
	return nil
}
func (s *StationAffiliationConfigConstraintImpl) Set(handler *action.Handler, newVal any) error {
	str, err := GetConfigValString(newVal)
	if err != nil {
		return err
	}
	if str == "" {
		return errors.New("station affiliation cannot be empty")
	}
	if len(str) > 50 {
		return errors.New("station affiliation should be less than 50 characters")
	}
	if err := handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), str); err != nil {
		return fmt.Errorf("failed to set station affiliation: %w", err)
	}
	return nil
}
func (s *StationAffiliationConfigConstraintImpl) Get(handler *action.Handler) (any, error) {
	val, _, _, err := handler.SettingsGet(s.GetNamespace(), s.GetKey())
	if err != nil {
		return nil, fmt.Errorf("failed to get station affiliation: %w", err)
	}
	if _, ok := val.(string); !ok {
		return nil, errors.New("string expected")
	}
	return val, nil
}
func (s *StationAffiliationConfigConstraintImpl) Restore(handler *action.Handler) error {
	if err := handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to reset station affiliation: %w", err)
	}
	return nil
}

type StationStationCodeConfigConstraintImpl struct{}

func (s *StationStationCodeConfigConstraintImpl) GetName() string             { return "Station Code" }
func (s *StationStationCodeConfigConstraintImpl) GetNamespace() string        { return STATION_NAMESPACE }
func (s *StationStationCodeConfigConstraintImpl) GetKey() string              { return "station_code" }
func (s *StationStationCodeConfigConstraintImpl) GetType() action.SettingType { return action.String }
func (s *StationStationCodeConfigConstraintImpl) IsRequired() bool            { return true }
func (s *StationStationCodeConfigConstraintImpl) GetVersion() int             { return 0 }
func (s *StationStationCodeConfigConstraintImpl) GetOptions() map[string]any  { return nil }
func (s *StationStationCodeConfigConstraintImpl) GetDefaultValue() any        { return "SHAKE" }
func (s *StationStationCodeConfigConstraintImpl) GetDescription() string {
	return "An unique 5-letter uppercase station code, used for services such as MiniSEED, Helicorder, SeedLink, etc."
}
func (s *StationStationCodeConfigConstraintImpl) Init(handler *action.Handler) error {
	if _, err := handler.SettingsInit(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to set default station code: %w", err)
	}
	return nil
}
func (s *StationStationCodeConfigConstraintImpl) Set(handler *action.Handler, newVal any) error {
	str, err := GetConfigValString(newVal)
	if err != nil {
		return err
	}
	if str == "" {
		return errors.New("station code cannot be empty")
	}
	if len(str) > 5 {
		return errors.New("station code should be less than 5 characters")
	}
	if err := handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), strings.ToUpper(str)); err != nil {
		return fmt.Errorf("failed to set station code: %w", err)
	}
	return nil
}
func (s *StationStationCodeConfigConstraintImpl) Get(handler *action.Handler) (any, error) {
	val, _, _, err := handler.SettingsGet(s.GetNamespace(), s.GetKey())
	if err != nil {
		return nil, fmt.Errorf("failed to get station code: %w", err)
	}
	if _, ok := val.(string); !ok {
		return nil, errors.New("string expected")
	}
	return val, nil
}
func (s *StationStationCodeConfigConstraintImpl) Restore(handler *action.Handler) error {
	if err := handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to reset station code: %w", err)
	}
	return nil
}

type StationNetworkCodeConfigConstraintImpl struct{}

func (s *StationNetworkCodeConfigConstraintImpl) GetName() string             { return "Network Code" }
func (s *StationNetworkCodeConfigConstraintImpl) GetNamespace() string        { return STATION_NAMESPACE }
func (s *StationNetworkCodeConfigConstraintImpl) GetKey() string              { return "network_code" }
func (s *StationNetworkCodeConfigConstraintImpl) GetType() action.SettingType { return action.String }
func (s *StationNetworkCodeConfigConstraintImpl) IsRequired() bool            { return true }
func (s *StationNetworkCodeConfigConstraintImpl) GetVersion() int             { return 0 }
func (s *StationNetworkCodeConfigConstraintImpl) GetOptions() map[string]any  { return nil }
func (s *StationNetworkCodeConfigConstraintImpl) GetDefaultValue() any        { return "AS" }
func (s *StationNetworkCodeConfigConstraintImpl) GetDescription() string {
	return "A unique 2-letter uppercase network code, used for services such as MiniSEED, Helicorder, SeedLink, etc."
}
func (s *StationNetworkCodeConfigConstraintImpl) Init(handler *action.Handler) error {
	if _, err := handler.SettingsInit(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to set default network code: %w", err)
	}
	return nil
}
func (s *StationNetworkCodeConfigConstraintImpl) Set(handler *action.Handler, newVal any) error {
	str, err := GetConfigValString(newVal)
	if err != nil {
		return err
	}
	if len(str) != 2 {
		return errors.New("network code should be 2 characters")
	}
	if err := handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), strings.ToUpper(str)); err != nil {
		return fmt.Errorf("failed to set network code: %w", err)
	}
	return nil
}
func (s *StationNetworkCodeConfigConstraintImpl) Get(handler *action.Handler) (any, error) {
	val, _, _, err := handler.SettingsGet(s.GetNamespace(), s.GetKey())
	if err != nil {
		return nil, fmt.Errorf("failed to get network code: %w", err)
	}
	if _, ok := val.(string); !ok {
		return nil, errors.New("failed to assert network code: string expected")
	}
	return val, nil
}
func (s *StationNetworkCodeConfigConstraintImpl) Restore(handler *action.Handler) error {
	if err := handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to reset network code: %w", err)
	}
	return nil
}

type StationLocationCodeConfigConstraintImpl struct{}

func (s *StationLocationCodeConfigConstraintImpl) GetName() string             { return "Location Code" }
func (s *StationLocationCodeConfigConstraintImpl) GetNamespace() string        { return STATION_NAMESPACE }
func (s *StationLocationCodeConfigConstraintImpl) GetKey() string              { return "location_code" }
func (s *StationLocationCodeConfigConstraintImpl) GetType() action.SettingType { return action.String }
func (s *StationLocationCodeConfigConstraintImpl) IsRequired() bool            { return true }
func (s *StationLocationCodeConfigConstraintImpl) GetVersion() int             { return 0 }
func (s *StationLocationCodeConfigConstraintImpl) GetOptions() map[string]any  { return nil }
func (s *StationLocationCodeConfigConstraintImpl) GetDefaultValue() any        { return "00" }
func (s *StationLocationCodeConfigConstraintImpl) GetDescription() string {
	return "A unique 2-letter uppercase location code, used for services such as MiniSEED, Helicorder, SeedLink, etc."
}
func (s *StationLocationCodeConfigConstraintImpl) Init(handler *action.Handler) error {
	if _, err := handler.SettingsInit(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to set default location code: %w", err)
	}
	return nil
}
func (s *StationLocationCodeConfigConstraintImpl) Set(handler *action.Handler, newVal any) error {
	str, err := GetConfigValString(newVal)
	if err != nil {
		return err
	}
	if len(str) != 2 {
		return errors.New("location code should be 2 characters")
	}
	if err := handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), strings.ToUpper(str)); err != nil {
		return fmt.Errorf("failed to set location code: %w", err)
	}
	return nil
}
func (s *StationLocationCodeConfigConstraintImpl) Get(handler *action.Handler) (any, error) {
	val, _, _, err := handler.SettingsGet(s.GetNamespace(), s.GetKey())
	if err != nil {
		return nil, fmt.Errorf("failed to get location code: %w", err)
	}
	if _, ok := val.(string); !ok {
		return nil, errors.New("failed to assert location code: string expected")
	}
	return val, nil
}
func (s *StationLocationCodeConfigConstraintImpl) Restore(handler *action.Handler) error {
	if err := handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to reset location code: %w", err)
	}
	return nil
}

type StationChannelCodesConfigConstraintImpl struct{}

func (s *StationChannelCodesConfigConstraintImpl) GetName() string      { return "Channel Codes" }
func (s *StationChannelCodesConfigConstraintImpl) GetNamespace() string { return STATION_NAMESPACE }
func (s *StationChannelCodesConfigConstraintImpl) GetKey() string       { return "channel_codes" }
func (s *StationChannelCodesConfigConstraintImpl) GetType() action.SettingType {
	return action.StringArray
}
func (s *StationChannelCodesConfigConstraintImpl) IsRequired() bool           { return true }
func (s *StationChannelCodesConfigConstraintImpl) GetOptions() map[string]any { return nil }
func (s *StationChannelCodesConfigConstraintImpl) GetVersion() int            { return 0 }
func (s *StationChannelCodesConfigConstraintImpl) GetDefaultValue() any {
	return []string{"EHZ", "EHE", "EHN", "ENZ", "ENE", "ENN", "CH7", "CH8"} // maximum 8 channels capacity in v3 protocol
}
func (s *StationChannelCodesConfigConstraintImpl) GetDescription() string {
	return "Code for each channel in uppercase, by default, the first three channels are set to \"EH*\", followed by three \"EN*\" sensors, software restart needed when channel codes are updated."
}
func (s *StationChannelCodesConfigConstraintImpl) Init(handler *action.Handler) error {
	channelCodes, ok := s.GetDefaultValue().([]string)
	if !ok {
		return errors.New("failed to assert channel codes")
	}
	if _, err := handler.SettingsInit(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), channelCodes); err != nil {
		return fmt.Errorf("failed to set default channel codes: %w", err)
	}
	return nil
}
func (s *StationChannelCodesConfigConstraintImpl) Set(handler *action.Handler, newVal any) error {
	channelCodes, err := GetConfigValStringArray(newVal)
	if err != nil {
		return err
	}
	var channelCodeStrArr []string
	for idx := 0; idx < len(channelCodes); idx++ {
		channelCode := channelCodes[idx]
		channelCode = fmt.Sprintf("%3s", strings.ToUpper(channelCode))
		if len(channelCode) > 3 {
			channelCode = channelCode[:3]
		}
		channelCodeStrArr = append(channelCodeStrArr, channelCode)
	}
	if err := handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), channelCodeStrArr); err != nil {
		return fmt.Errorf("failed to set channel codes: %w", err)
	}
	return nil
}
func (s *StationChannelCodesConfigConstraintImpl) Get(handler *action.Handler) (any, error) {
	val, _, _, err := handler.SettingsGet(s.GetNamespace(), s.GetKey())
	if err != nil {
		return nil, fmt.Errorf("failed to get channel codes: %w", err)
	}
	if _, ok := val.([]string); !ok {
		return nil, errors.New("failed to assert channel codes: string array expected")
	}
	return val, nil
}
func (s *StationChannelCodesConfigConstraintImpl) Restore(handler *action.Handler) error {
	if err := handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to reset channel codes: %w", err)
	}
	return nil
}

func NewStationConstraints() []IConstraint {
	return []IConstraint{
		&StationNameConfigConstraintImpl{},
		&StationDescriptionConfigConstraintImpl{},
		&StationCountryConfigConstraintImpl{},
		&StationPlaceConfigConstraintImpl{},
		&StationAffiliationConfigConstraintImpl{},
		&StationStationCodeConfigConstraintImpl{},
		&StationNetworkCodeConfigConstraintImpl{},
		&StationLocationCodeConfigConstraintImpl{},
		&StationChannelCodesConfigConstraintImpl{},
	}
}
