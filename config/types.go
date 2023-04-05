package config

type Args struct {
	Path    string
	Version bool
}

type Config struct {
	// 测站配置
	Station struct {
		UUID      string  `json:"uuid"`
		Name      string  `json:"name"`
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
		Altitude  float64 `json:"altitude"`
	} `json:"station_settings"`
	// 下位机配置
	Geophone struct {
		Device      string  `json:"device"`
		Baud        int     `json:"baud"`
		Sensitivity float64 `json:"sensitivity"`
	} `json:"geophone_settings"`
	// NTP 配置
	NTPClient struct {
		Host     string `json:"host"`
		Port     int    `json:"port"`
		Timeout  int    `json:"timeout"`
		Interval int    `json:"interval"`
	} `json:"ntpclient_settings"`
	// 数据推送配置
	Push struct {
		Host string `json:"host"`
		Port int    `json:"port"`
		Path string `json:"path"`
		TLS  bool   `json:"tls"`
	} `json:"push_settings"`
	// Web 配置
	Web struct {
		Host string `json:"host"`
		Port int    `json:"port"`
	} `json:"web_settings"`
}
