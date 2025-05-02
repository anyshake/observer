package quakesense

import (
	"fmt"

	"github.com/anyshake/observer/config"
)

func (s *QuakeSenseServiceImpl) Init() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, con := range s.GetConfigConstraint() {
		if err := con.Init(s.actionHandler); err != nil {
			return fmt.Errorf("failed to initlize config constraint for service %s, namespace %s, key %s: %w", ID, con.GetNamespace(), con.GetKey(), err)
		}
	}

	mqttBroker, err := (&quakeSenseConfigMqttBrokerImpl{}).Get(s.actionHandler)
	if err != nil {
		return fmt.Errorf("failed to get quakeSense MQTT broker: %w", err)
	}
	s.mqttBroker = mqttBroker.(string)

	mqttTopic, err := (&quakeSenseConfigMqttTopicImpl{}).Get(s.actionHandler)
	if err != nil {
		return fmt.Errorf("failed to get quakeSense MQTT topic: %w", err)
	}
	s.mqttTopic = mqttTopic.(string)

	mqttUsername, err := (&quakeSenseConfigMqttUsernameImpl{}).Get(s.actionHandler)
	if err != nil {
		return fmt.Errorf("failed to get quakeSense MQTT username: %w", err)
	}
	s.mqttUsername = mqttUsername.(string)

	mqttPassword, err := (&quakeSenseConfigMqttPasswordImpl{}).Get(s.actionHandler)
	if err != nil {
		return fmt.Errorf("failed to get quakeSense MQTT password: %w", err)
	}
	s.mqttPassword = mqttPassword.(string)

	stationName, err := (&config.StationNameConfigConstraintImpl{}).Get(s.actionHandler)
	if err != nil {
		return fmt.Errorf("failed to get quakeSense station name: %w", err)
	}
	s.stationName = stationName.(string)

	stationDescription, err := (&config.StationDescriptionConfigConstraintImpl{}).Get(s.actionHandler)
	if err != nil {
		return fmt.Errorf("failed to get quakeSense station description: %w", err)
	}
	s.stationDescription = stationDescription.(string)

	stationPlace, err := (&config.StationPlaceConfigConstraintImpl{}).Get(s.actionHandler)
	if err != nil {
		return fmt.Errorf("failed to get quakeSense station place: %w", err)
	}
	s.stationPlace = stationPlace.(string)

	stationCountry, err := (&config.StationCountryConfigConstraintImpl{}).Get(s.actionHandler)
	if err != nil {
		return fmt.Errorf("failed to get quakeSense station country: %w", err)
	}
	s.stationCountry = stationCountry.(string)

	stationAffiliation, err := (&config.StationAffiliationConfigConstraintImpl{}).Get(s.actionHandler)
	if err != nil {
		return fmt.Errorf("failed to get quakeSense station affiliation: %w", err)
	}
	s.stationAffiliation = stationAffiliation.(string)

	stationCode, err := (&config.StationStationCodeConfigConstraintImpl{}).Get(s.actionHandler)
	if err != nil {
		return fmt.Errorf("failed to get quakeSense station code: %w", err)
	}
	s.stationCode = stationCode.(string)

	networkCode, err := (&config.StationNetworkCodeConfigConstraintImpl{}).Get(s.actionHandler)
	if err != nil {
		return fmt.Errorf("failed to get quakeSense network code: %w", err)
	}
	s.networkCode = networkCode.(string)

	locationCode, err := (&config.StationLocationCodeConfigConstraintImpl{}).Get(s.actionHandler)
	if err != nil {
		return fmt.Errorf("failed to get quakeSense location code: %w", err)
	}
	s.locationCode = locationCode.(string)

	monitorChannel, err := (&quakeSenseConfigMonitorChannelImpl{}).Get(s.actionHandler)
	if err != nil {
		return fmt.Errorf("failed to get quakeSense monitor channel: %w", err)
	}
	s.monitorChannel = monitorChannel.(string)

	throttleSeconds, err := (&quakeSenseConfigThrottleImpl{}).Get(s.actionHandler)
	if err != nil {
		return fmt.Errorf("failed to get quakeSense throttleSeconds: %w", err)
	}
	s.throttleSeconds = int(throttleSeconds.(int64))

	triggerMethod, err := (&quakeSenseConfigTriggerMethodImpl{}).Get(s.actionHandler)
	if err != nil {
		return fmt.Errorf("failed to get quakeSense trigger method: %w", err)
	}
	s.triggerMethod = triggerMethod.(string)

	staWindow, err := (&quakeSenseConfigStaWindowImpl{}).Get(s.actionHandler)
	if err != nil {
		return fmt.Errorf("failed to get quakeSense STA window: %w", err)
	}
	s.staWindow = staWindow.(float64)

	ltaWindow, err := (&quakeSenseConfigLtaWindowImpl{}).Get(s.actionHandler)
	if err != nil {
		return fmt.Errorf("failed to get quakeSense LTA window: %w", err)
	}
	s.ltaWindow = ltaWindow.(float64)

	if s.ltaWindow <= s.staWindow {
		return fmt.Errorf("LTA window cannot be smaller than STA window")
	}

	trigOn, err := (&quakeSenseConfigTrigOnImpl{}).Get(s.actionHandler)
	if err != nil {
		return fmt.Errorf("failed to get quakeSense trigger on: %w", err)
	}
	s.trigOn = trigOn.(float64)

	trigOff, err := (&quakeSenseConfigTrigOffImpl{}).Get(s.actionHandler)
	if err != nil {
		return fmt.Errorf("failed to get quakeSense trigger off: %w", err)
	}
	s.trigOff = trigOff.(float64)

	filterType, err := (&quakeSenseConfigFilterTypeImpl{}).Get(s.actionHandler)
	if err != nil {
		return fmt.Errorf("failed to get quakeSense filter type: %w", err)
	}
	s.filterType = filterType.(string)

	maxFreq, err := (&quakeSenseConfigMaxFreqImpl{}).Get(s.actionHandler)
	if err != nil {
		return fmt.Errorf("failed to get quakeSense max frequency: %w", err)
	}
	s.maxFreq = maxFreq.(float64)

	minFreq, err := (&quakeSenseConfigMinFreqImpl{}).Get(s.actionHandler)
	if err != nil {
		return fmt.Errorf("failed to get quakeSense min frequency: %w", err)
	}
	s.minFreq = minFreq.(float64)

	return nil
}
