package config

import "github.com/spf13/viper"

type Config struct {
	HTTP_PORT         string `mapstructure:"HTTP_PORT"`
	POSTGRES_USER     string `mapstructure:"POSTGRES_USER"`
	POSTGRES_PASSWORD string `mapstructure:"POSTGRES_PASSWORD"`
	POSTGRES_HOSTNAME string `mapstructure:"POSTGRES_HOSTNAME"`
	POSTGRES_PORT     string `mapstructure:"POSTGRES_PORT"`
	POSTGRES_DB       string `mapstructure:"POSTGRES_DB"`
	POSTGRES_SSL      string `mapstructure:"POSTGRES_SSL"`
	REDIS_ADDRESS     string `mapstructure:"REDIS_ADDRESS"`
}

func LoadConfig() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.ReadInConfig()
	cfg := Config{}
	err := viper.Unmarshal(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
