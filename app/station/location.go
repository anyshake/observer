package station

import "com.geophone.observer/features/collector"

func GetLocation(message collector.Message) Location {
	return Location{
		Latitude:  message.Acceleration[0].Latitude,
		Longitude: message.Acceleration[0].Longitude,
		Altitude:  message.Acceleration[0].Altitude,
	}
}
