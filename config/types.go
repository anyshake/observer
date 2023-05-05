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
		Device      string `json:"device"`
		Baud        int    `json:"baud"`
		Sensitivity struct {
			Vertical   float64 `json:"vertical"`
			EastWest   float64 `json:"east_west"`
			NorthSouth float64 `json:"north_south"`
		} `json:"sensitivity"`
	} `json:"geophone_settings"`
	// NTP 配置
	NTPClient struct {
		Host     string `json:"host"`
		Port     int    `json:"port"`
		Timeout  int    `json:"timeout"`
		Interval int    `json:"interval"`
	} `json:"ntpclient_settings"`
	// 数据收集配置
	Collector struct {
		Enable bool   `json:"enable"`
		Host   string `json:"host"`
		Port   int    `json:"port"`
		Path   string `json:"path"`
		TLS    bool   `json:"tls"`
	} `json:"collector_settings"`
	// 数据存档设定
	Archiver struct {
		Enable   bool   `json:"enable"`
		Host     string `json:"host"`
		Port     int    `json:"port"`
		Username string `json:"username"`
		Password string `json:"password"`
		Database string `json:"database"`
	} `json:"archiver_settings"`
	// HTTP 服务器配置
	Server struct {
		Host string `json:"host"`
		Port int    `json:"port"`
	} `json:"server_settings"`
}
