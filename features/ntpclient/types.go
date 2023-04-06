package ntpclient

type NTP struct {
	Offset float64
}

type NTPOptions struct {
	NTP             *NTP
	Port            int
	Timeout         int
	OnDataCallback  func(*NTP)
	OnErrorCallback func(error)
}
