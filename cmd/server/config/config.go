package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	HTTP_PORT         string `mapstructure:"HTTP_PORT"`
	POSTGRES_USER     string `mapstructure:"POSTGRES_USER"`
	POSTGRES_PASSWORD string `mapstructure:"POSTGRES_PASSWORD"`
	POSTGRES_HOSTNAME string `mapstructure:"POSTGRES_HOSTNAME"`
	POSTGRES_PORT     string `mapstructure:"POSTGRES_PORT"`
	POSTGRES_DB       string `mapstructure:"POSTGRES_DB"`
	POSTGRES_SSL      string `mapstructure:"POSTGRES_SSL"`
	REDIS_HOSTNAME    string `mapstructure:"REDIS_HOSTNAME"`
	REDIS_PORT        string `mapstructure:"REDIS_PORT"`
}

func LoadConfig() (*Config, error) {
	var config Config

	// Load .env file manually only if vars aren't already in environment
	if os.Getenv("HTTP_PORT") == "" {
		_ = godotenv.Load(".env")
	}

	// Set Viper to read from ENV
	viper.AutomaticEnv()

	// Bind to struct
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("unable to decode into struct: %w", err)
	}

	return &config, nil
}
