package config

type station struct {
	UUID      string  `json:"uuid"`
	Name      string  `json:"name"`
	Station   string  `json:"station"`
	Network   string  `json:"network"`
	Location  string  `json:"location"`
	Country   string  `json:"country"`
	Region    string  `json:"region"`
	City      string  `json:"city"`
	Owner     string  `json:"owner"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Elevation float64 `json:"elevation"`
}

type serial struct {
	Device string `json:"device"`
	Baud   int    `json:"baud"`
	Packet int    `json:"packet"`
}

type adc struct {
	FullScale  float64 `json:"fullscale"`
	Resolution int     `json:"resolution"`
}

type geophone struct {
	Frequency   float64 `json:"frequency"`
	Sensitivity float64 `json:"sensitivity"`
}

type ntpclient struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Timeout  int    `json:"timeout"`
	Interval int    `json:"interval"`
}

type archiver struct {
	Engine    string `json:"engine"`
	Enable    bool   `json:"enable"`
	Host      string `json:"host"`
	Port      int    `json:"port"`
	LifeCycle int    `json:"lifecycle"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	Database  string `json:"database"`
}

type server struct {
	Host  string `json:"host"`
	Port  int    `json:"port"`
	CORS  bool   `json:"cors"`
	Debug bool   `json:"debug"`
	Rate  int    `json:"rate"`
}

type miniseed struct {
	Enable    bool   `json:"enable"`
	Path      string `json:"path"`
	LifeCycle int    `json:"lifecycle"`
}

type seedlink struct {
	Enable   bool   `json:"enable"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Duration int    `json:"duration"`
}

type Conf struct {
	Station   station   `json:"station_settings"`
	Serial    serial    `json:"serial_settings"`
	ADC       adc       `json:"adc_settings"`
	Geophone  geophone  `json:"geophone_settings"`
	NTPClient ntpclient `json:"ntpclient_settings"`
	Archiver  archiver  `json:"archiver_settings"`
	Server    server    `json:"server_settings"`
	MiniSEED  miniseed  `json:"miniseed_settings"`
	SeedLink  seedlink  `json:"seedlink_settings"`
}

type Args struct {
	Path    string // Path to config file
	Version bool   // Show version information
}
