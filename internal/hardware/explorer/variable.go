package explorer

import (
	"errors"
	"math"
)

func (v *DeviceVariable) Reset() {
	v.mu.Lock()
	defer v.mu.Unlock()

	v.deviceId = nil
	v.latitude = nil
	v.longitude = nil
	v.elevation = nil
	v.temperature = nil
}

func (v *DeviceVariable) SetDeviceId(deviceId *uint32) {
	v.mu.Lock()
	defer v.mu.Unlock()

	if deviceId == nil {
		v.deviceId = nil
		return
	}

	devId := *deviceId
	if devId == 0x7FFFFFFF {
		devId = math.MaxUint32
	}

	if v.deviceId == nil {
		v.deviceId = &devId
		return
	}

	if *v.deviceId != devId {
		v.deviceId = &devId
	}
}

func (v *DeviceVariable) GetDeviceId() (uint32, error) {
	v.mu.Lock()
	defer v.mu.Unlock()

	if v.deviceId == nil {
		return 0, errors.New("device id is not set")
	}

	return *v.deviceId, nil
}

func (v *DeviceVariable) SetLatitude(latitude *float64) {
	v.mu.Lock()
	defer v.mu.Unlock()

	if latitude == nil {
		v.latitude = nil
		return
	}

	if v.latitude == nil {
		v.latitude = latitude
		return
	}

	if *v.latitude != *latitude {
		v.latitude = latitude
	}
}

func (v *DeviceVariable) GetLatitude() (float64, error) {
	v.mu.Lock()
	defer v.mu.Unlock()

	if v.latitude == nil {
		return 0, errors.New("latitude is not set")
	}

	return *v.latitude, nil
}

func (v *DeviceVariable) SetLongitude(longitude *float64) {
	v.mu.Lock()
	defer v.mu.Unlock()

	if longitude == nil {
		v.longitude = nil
		return
	}

	if v.longitude == nil {
		v.longitude = longitude
		return
	}

	if *v.longitude != *longitude {
		v.longitude = longitude
	}
}

func (v *DeviceVariable) GetLongitude() (float64, error) {
	v.mu.Lock()
	defer v.mu.Unlock()

	if v.longitude == nil {
		return 0, errors.New("longitude is not set")
	}

	return *v.longitude, nil
}

func (v *DeviceVariable) SetElevation(elevation *float64) {
	v.mu.Lock()
	defer v.mu.Unlock()

	if elevation == nil {
		v.elevation = nil
		return
	}

	if v.elevation == nil {
		v.elevation = elevation
		return
	}

	if *v.elevation != *elevation {
		v.elevation = elevation
	}
}

func (v *DeviceVariable) GetElevation() (float64, error) {
	v.mu.Lock()
	defer v.mu.Unlock()

	if v.elevation == nil {
		return 0, errors.New("elevation is not set")
	}

	return *v.elevation, nil
}

func (v *DeviceVariable) SetTemperature(temperature *float64) {
	v.mu.Lock()
	defer v.mu.Unlock()

	if temperature == nil {
		v.temperature = nil
		return
	}

	if v.temperature == nil {
		v.temperature = temperature
		return
	}

	if *v.temperature != *temperature {
		v.temperature = temperature
	}
}

func (v *DeviceVariable) GetTemperature() (float64, error) {
	v.mu.Lock()
	defer v.mu.Unlock()

	if v.temperature == nil {
		return 0, errors.New("temperature is not set")
	}

	return *v.temperature, nil
}
