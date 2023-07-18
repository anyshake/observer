package config

type Args struct {
	Path    string // Path to config file
	Version bool   // Show version information
}

type Conf struct {
	Station struct {
		UUID      string  `json:"uuid"`
		Name      string  `json:"name"`
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
		Altitude  float64 `json:"altitude"`
	} `json:"station_settings"`
	Serial struct {
		Device string `json:"device"`
		Baud   int    `json:"baud"`
	} `json:"serial_settings"`
	ADC struct {
		FullScale  float64 `json:"fullscale"`
		Resolution int     `json:"resolution"`
	} `json:"adc_settings"`
	Geophone struct {
		EHZ float64 `json:"ehz"`
		EHE float64 `json:"ehe"`
		EHN float64 `json:"ehn"`
	} `json:"geophone_settings"`
	NTPClient struct {
		Host     string `json:"host"`
		Port     int    `json:"port"`
		Timeout  int    `json:"timeout"`
		Interval int    `json:"interval"`
	} `json:"ntpclient_settings"`
	Archiver struct {
		Enable   bool   `json:"enable"`
		Host     string `json:"host"`
		Port     int    `json:"port"`
		Username string `json:"username"`
		Password string `json:"password"`
		Database string `json:"database"`
	} `json:"archiver_settings"`
	Server struct {
		Host string `json:"host"`
		Port int    `json:"port"`
	} `json:"server_settings"`
}
