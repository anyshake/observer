package config

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type location struct {
	Latitude  float64 `mapstructure:"latitude" validate:"lte=90,gte=-90"`
	Longitude float64 `mapstructure:"longitude" validate:"lte=180,gte=-180"`
	Elevation float64 `mapstructure:"elevation"`
}

type hardware struct {
	Endpoint string `mapstructure:"endpoint" validate:"required"`
	Protocol string `mapstructure:"protocol" validate:"required"`
	Model    string `mapstructure:"model"`
	Timeout  int    `mapstructure:"timeout" validate:"gte=0"`
}

type ntpClient struct {
	// Endpoint is a deprecated field since v4.2.0, use Pool instead
	Endpoint string   `mapstructure:"endpoint"`
	Timeout  int      `mapstructure:"timeout" validate:"gte=0"`
	Retry    int      `mapstructure:"retry" validate:"gte=0"`
	Pool     []string `mapstructure:"pool"`
}

type database struct {
	Endpoint string `mapstructure:"endpoint" validate:"required"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Prefix   string `mapstructure:"prefix"`
	Database string `mapstructure:"database" validate:"required"`
	Timeout  int    `mapstructure:"timeout" validate:"gte=0"`
}

type server struct {
	Listen string `mapstructure:"listen" validate:"required"`
	Debug  bool   `mapstructure:"debug"`
	CORS   bool   `mapstructure:"cors"`
}

type logger struct {
	Level     string `mapstructure:"level" validate:"oneof=info warn error"`
	LifeCycle int    `mapstructure:"lifecycle" validate:"gte=0"`
	Rotation  int    `mapstructure:"rotation" validate:"gte=0"`
	Size      int    `mapstructure:"size" validate:"gte=0"`
	Path      string `mapstructure:"path"`
}

type BaseConfig struct {
	Location  location  `mapstructure:"location"`
	Hardware  hardware  `mapstructure:"hardware"`
	NtpClient ntpClient `mapstructure:"ntpclient"`
	Database  database  `mapstructure:"database"`
	Server    server    `mapstructure:"server"`
	Logger    logger    `mapstructure:"logger"`
}

func (c *BaseConfig) Parse(configPath, configType string) error {
	viper.SetConfigFile(configPath)
	viper.SetConfigType(configType)

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read config: %w", err)
	}

	if err := viper.Unmarshal(c); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	validate := validator.New()
	if err := validate.Struct(c); err != nil {
		return fmt.Errorf("failed to validate config: %w", err)
	}

	return nil
}
