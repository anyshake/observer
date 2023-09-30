package config

type station struct {
	UUID      string  `json:"uuid"`
	Name      string  `json:"name"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Altitude  float64 `json:"altitude"`
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

type channel struct {
	Compensation bool    `json:"compensation"`
	Damping      float64 `json:"damping"`
	Frequency    float64 `json:"frequency"`
	Sensitivity  float64 `json:"sensitivity"`
}

type geophone struct {
	EHZ channel `json:"ehz"`
	EHE channel `json:"ehe"`
	EHN channel `json:"ehn"`
}

type ntpclient struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Timeout  int    `json:"timeout"`
	Interval int    `json:"interval"`
}

type archiver struct {
	Enable   bool   `json:"enable"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Database string `json:"database"`
}

type server struct {
	Host  string `json:"host"`
	Port  int    `json:"port"`
	CORS  bool   `json:"cors"`
	Debug bool   `json:"debug"`
}

type miniseed struct {
	Enable  bool   `json:"enable"`
	Path    string `json:"path"`
	Station string `json:"station"`
	Network string `json:"network"`
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
}

type Args struct {
	Path    string // Path to config file
	Version bool   // Show version information
}
