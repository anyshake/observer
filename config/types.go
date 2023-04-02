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
	// 日志配置
	Log struct {
		Level string `json:"level"`
		Path  string `json:"path"`
	} `json:"log_settings"`
	// ADC 配置
	ADC struct {
		Depth int     `json:"depth"`
		VRef  float64 `json:"vref"`
		Gain  float64 `json:"gain"`
	} `json:"adc_settings"`
	// 下位机配置
	Geophone struct {
		Device      string  `json:"device"`
		Baud        int     `json:"baud"`
		Sensitivity float64 `json:"sensitivity"`
	} `json:"geophone_settings"`
	// GNSS 配置
	GNSS struct {
		Host string `json:"host"`
		Port int    `json:"port"`
	} `json:"gnss_settings"`
	// 采集器配置
	Collector struct {
		Host string `json:"host"`
		Port int    `json:"port"`
		Path string `json:"path"`
	} `json:"collector_settings"`
	// Web 配置
	Web struct {
		Host string `json:"host"`
		Port int    `json:"port"`
	} `json:"web_settings"`
}
