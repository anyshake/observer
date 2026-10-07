package explorer

import (
	"math"
	"testing"
	"time"
)

func TestDeviceConfigAccessors(t *testing.T) {
	t.Parallel()

	var config DeviceConfig
	if config.GetPacketInterval() != 0 || config.GetSampleRate() != 0 || config.GetGnssAvailability() || config.GetModel() != "" || config.GetProtocol() != "" || config.GetChannelCodes() != nil {
		t.Fatalf("zero config = %#v", config)
	}

	config.SetPacketInterval(200 * time.Millisecond)
	config.SetSampleRate(100)
	config.SetGnssAvailability(true)
	config.SetModel("E-C111G")
	config.SetProtocol("v3")
	config.SetChannelCodes([]string{"EHZ", "EHE", "EHN"})

	if config.GetPacketInterval() != 200*time.Millisecond || config.GetSampleRate() != 100 || !config.GetGnssAvailability() {
		t.Fatalf("unexpected numeric config: %#v", config)
	}
	if config.GetModel() != "E-C111G" || config.GetProtocol() != "v3" {
		t.Fatalf("unexpected identity config: %#v", config)
	}
	if got := config.GetChannelCodes(); len(got) != 3 || got[0] != "EHZ" || got[2] != "EHN" {
		t.Fatalf("channel codes = %#v", got)
	}

	config.SetGnssAvailability(false)
	config.SetChannelCodes(nil)
	if config.GetGnssAvailability() || config.GetChannelCodes() != nil {
		t.Fatal("config did not accept cleared values")
	}
}

func TestDeviceStatusAccessors(t *testing.T) {
	t.Parallel()

	var status DeviceStatus
	if !status.GetStartedAt().IsZero() || !status.GetUpdatedAt().IsZero() || status.GetFrames() != 0 || status.GetErrors() != 0 || status.GetMessages() != 0 {
		t.Fatalf("zero status = %#v", status)
	}

	started := time.Unix(1700000000, 0)
	updated := started.Add(time.Second)
	status.SetStartedAt(started)
	status.SetUpdatedAt(updated)
	status.IncrementFrames()
	status.IncrementFrames()
	status.IncrementErrors()
	status.IncrementMessages()
	status.IncrementMessages()
	status.IncrementMessages()

	if !status.GetStartedAt().Equal(started) || !status.GetUpdatedAt().Equal(updated) {
		t.Fatalf("timestamps = %v %v", status.GetStartedAt(), status.GetUpdatedAt())
	}
	if status.GetFrames() != 2 || status.GetErrors() != 1 || status.GetMessages() != 3 {
		t.Fatalf("counters = %d frames, %d errors, %d messages", status.GetFrames(), status.GetErrors(), status.GetMessages())
	}
}

func TestDeviceVariableAccessors(t *testing.T) {
	t.Parallel()

	var variable DeviceVariable
	if _, err := variable.GetDeviceId(); err == nil {
		t.Fatal("expected unset device id")
	}
	if _, err := variable.GetLatitude(false); err == nil {
		t.Fatal("expected unset latitude")
	}
	if _, err := variable.GetLongitude(true); err == nil {
		t.Fatal("expected unset longitude")
	}
	if _, err := variable.GetElevation(); err == nil {
		t.Fatal("expected unset elevation")
	}
	if _, err := variable.GetTemperature(); err == nil {
		t.Fatal("expected unset temperature")
	}

	variable.SetDeviceId(nil)
	variable.SetLatitude(nil)
	variable.SetLongitude(nil)
	variable.SetElevation(nil)
	variable.SetTemperature(nil)
	if _, err := variable.GetDeviceId(); err == nil {
		t.Fatal("nil setters left a device id")
	}

	sentinel := uint32(0x7FFFFFFF)
	variable.SetDeviceId(&sentinel)
	gotID, err := variable.GetDeviceId()
	if err != nil || gotID != math.MaxUint32 {
		t.Fatalf("saturated device id = %08X, %v", gotID, err)
	}
	variable.SetDeviceId(&sentinel)
	if gotID, err = variable.GetDeviceId(); err != nil || gotID != math.MaxUint32 {
		t.Fatalf("repeat saturated device id = %08X, %v", gotID, err)
	}

	deviceID := uint32(0x012F81AC)
	variable.SetDeviceId(&deviceID)
	variable.SetDeviceId(&deviceID)
	if gotID, err = variable.GetDeviceId(); err != nil || gotID != deviceID {
		t.Fatalf("device id = %08X, %v", gotID, err)
	}

	latitude := 25.126
	longitude := 121.234
	elevation := 42.5
	temperature := 18.25
	variable.SetLatitude(&latitude)
	variable.SetLatitude(&latitude)
	variable.SetLongitude(&longitude)
	variable.SetLongitude(&longitude)
	variable.SetElevation(&elevation)
	variable.SetElevation(&elevation)
	variable.SetTemperature(&temperature)
	variable.SetTemperature(&temperature)

	gotLat, err := variable.GetLatitude(false)
	if err != nil || gotLat != latitude {
		t.Fatalf("latitude = %v, %v", gotLat, err)
	}
	gotLat, err = variable.GetLatitude(true)
	if err != nil || gotLat != 25.13 {
		t.Fatalf("fuzzy latitude = %v, %v", gotLat, err)
	}
	gotLon, err := variable.GetLongitude(false)
	if err != nil || gotLon != longitude {
		t.Fatalf("longitude = %v, %v", gotLon, err)
	}
	gotLon, err = variable.GetLongitude(true)
	if err != nil || gotLon != 121.23 {
		t.Fatalf("fuzzy longitude = %v, %v", gotLon, err)
	}
	if gotElv, err := variable.GetElevation(); err != nil || gotElv != elevation {
		t.Fatalf("elevation = %v, %v", gotElv, err)
	}
	if gotTemp, err := variable.GetTemperature(); err != nil || gotTemp != temperature {
		t.Fatalf("temperature = %v, %v", gotTemp, err)
	}

	replacement := uint32(0x10)
	nextLatitude := latitude + 1
	nextLongitude := longitude + 1
	nextElevation := elevation + 1
	nextTemperature := temperature + 1
	variable.SetDeviceId(&replacement)
	variable.SetLatitude(&nextLatitude)
	variable.SetLongitude(&nextLongitude)
	variable.SetElevation(&nextElevation)
	variable.SetTemperature(&nextTemperature)
	if gotID, err = variable.GetDeviceId(); err != nil || gotID != replacement {
		t.Fatalf("replaced device id = %08X, %v", gotID, err)
	}
	if gotLat, err = variable.GetLatitude(false); err != nil || gotLat != nextLatitude {
		t.Fatalf("replaced latitude = %v, %v", gotLat, err)
	}

	variable.Reset()
	if _, err = variable.GetDeviceId(); err == nil {
		t.Fatal("reset left the device id set")
	}
	if _, err = variable.GetLatitude(false); err == nil {
		t.Fatal("reset left latitude set")
	}
	if _, err = variable.GetLongitude(false); err == nil {
		t.Fatal("reset left longitude set")
	}
	if _, err = variable.GetElevation(); err == nil {
		t.Fatal("reset left elevation set")
	}
	if _, err = variable.GetTemperature(); err == nil {
		t.Fatal("reset left temperature set")
	}
}
