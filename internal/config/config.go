package config

import (
	"errors"
	"fmt"

	"github.com/spf13/viper"
)

// Config holds all application configuration settings.
// Configuration values can be loaded from config/default.yaml file
// or overridden by environment variables.
type Config struct {
	// AppEnv specifies the app env
	AppEnv string `mapstructure:"APP_ENV"`

	// HTTPPort is the port number on which the HTTP server will listen.
	HTTPPort string `mapstructure:"HTTP_PORT"`

	// DBHost is the hostname or IP address of the database server.
	DBHost string `mapstructure:"DB_HOST"`

	// DBPort is the port number of the db server.
	DBPort int `mapstructure:"DB_PORT"`

	// DBUser is the database username for authentication.
	DBUser string `mapstructure:"DB_USER"`

	// DBPassword is the database password for authentication.
	DBPassword string `mapstructure:"DB_PASSWORD"`

	// DBName is the name of the database to connect to.
	DBName string `mapstructure:"DB_NAME"`
}

// Load reads configuration from file or environment variables.
func Load() (*Config, error) {
	v := viper.New()

	v.SetDefault("APP_ENV", "development")
	v.SetDefault("HTTP_PORT", "8080")
	v.SetDefault("DB_PORT", 5432)

	// load from config/default.yaml
	v.AddConfigPath("config")
	v.SetConfigName("default")
	v.SetConfigType("yaml")

	if err := v.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if !errors.As(err, &configFileNotFoundError) {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	// allows env var to override config file
	v.AutomaticEnv()

	// bind to Struct
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unable to decode into struct: %w", err)
	}

	return &cfg, nil
}
