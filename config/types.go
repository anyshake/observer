package config

type Station struct {
	Name    string `json:"name"`
	Owner   string `json:"owner"`
	Region  string `json:"region"`
	Country string `json:"country"`
	City    string `json:"city"`
}

type location struct {
	Latitude  float64 `json:"latitude" validate:"latitude"`
	Longitude float64 `json:"longitude" validate:"longitude"`
	Elevation float64 `json:"elevation" validate:"min=0"`
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
	Station  string `json:"station" validate:"max=5"`
	Network  string `json:"network" validate:"max=2"`
	Location string `json:"location" validate:"max=2"`
	Channel  string `json:"channel" validate:"max=3"`
}

type ntpclient struct {
	Host    string `json:"host"`
	Port    int    `json:"port" validate:"min=1,max=65535"`
	Timeout int    `json:"timeout" validate:"gte=0"`
	Retry   int    `json:"retry" validate:"gte=0"`
}

type database struct {
	Engine    string `json:"engine"`
	Host      string `json:"host"`
	Port      int    `json:"port" validate:"min=0,max=65535"`
	LifeCycle int    `json:"lifecycle" validate:"gte=0"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	Database  string `json:"database"`
}

type server struct {
	Host     string `json:"host"`
	Port     int    `json:"port" validate:"min=1,max=65535"`
	CORS     bool   `json:"cors"`
	Debug    bool   `json:"debug"`
	Restrict bool   `json:"restrict"`
	Rate     int    `json:"rate" validate:"gte=0"`
}

type logger struct {
	Level     string `json:"level" validate:"oneof=info warn error"`
	LifeCycle int    `json:"lifecycle" validate:"gte=0"`
	Rotation  int    `json:"rotation" validate:"gte=0"`
	Size      int    `json:"size" validate:"gte=0"`
	Dump      string `json:"dump"`
}

type services map[string]any

type Config struct {
	Station   Station   `json:"station_settings"`
	Location  location  `json:"location_settings"`
	Explorer  explorer  `json:"explorer_settings"`
	Sensor    Sensor    `json:"sensor_settings"`
	Stream    Stream    `json:"stream_settings"`
	NtpClient ntpclient `json:"ntpclient_settings"`
	Database  database  `json:"database_settings"`
	Server    server    `json:"server_settings"`
	Logger    logger    `json:"logger_settings"`
	Services  services  `json:"services_settings"`
}
