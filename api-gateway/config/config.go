package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Consul  Consul  `mapstructure:"consul"`
	Gateway Gateway `mapstructure:"gateway"`
}

type Consul struct {
	Address string `mapstructure:"address"`
}

type Gateway struct {
	Port string `mapstructure:"port"`
}

func LoadConfig() (Config, error) {
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	var cfg Config

	if err := viper.ReadInConfig(); err != nil {
		return cfg, err
	}

	return cfg, viper.Unmarshal(&cfg)
}
