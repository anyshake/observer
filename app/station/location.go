package station

import "com.geophone.observer/features/collector"

func GetLocation(message collector.Message) Location {
	return Location{
		Latitude:  message.Latitude,
		Longitude: message.Longitude,
		Altitude:  message.Altitude,
	}
}
