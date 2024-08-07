package config

type Station struct {
	Name    string `json:"name"`
	Owner   string `json:"owner"`
	Region  string `json:"region"`
	Country string `json:"country"`
	City    string `json:"city"`
}

type location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Elevation float64 `json:"elevation"`
}

type explorer struct {
	Legacy bool   `json:"legacy"`
	DSN    string `json:"dsn"`
	Engine string `json:"engine"`
}

type Sensor struct {
	Frequency   float64 `json:"frequency"`
	Sensitivity float64 `json:"sensitivity"`
	Velocity    bool    `json:"velocity"`
	Vref        float64 `json:"vref"`
	FullScale   float64 `json:"fullscale"`
	Resolution  int     `json:"resolution"`
}

type Stream struct {
	Station  string `json:"station"`
	Network  string `json:"network"`
	Location string `json:"location"`
	Channel  string `json:"channel"`
}

type ntpclient struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

type database struct {
	Engine    string `json:"engine"`
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

type logger struct {
	Level string `json:"level"`
	Dump  string `json:"dump"`
}

type Config struct {
	Station   Station        `json:"station_settings"`
	Location  location       `json:"location_settings"`
	Explorer  explorer       `json:"explorer_settings"`
	Sensor    Sensor         `json:"sensor_settings"`
	Stream    Stream         `json:"stream_settings"`
	NtpClient ntpclient      `json:"ntpclient_settings"`
	Database  database       `json:"database_settings"`
	Server    server         `json:"server_settings"`
	Logger    logger         `json:"logger_settings"`
	Services  map[string]any `json:"services_settings"`
}
