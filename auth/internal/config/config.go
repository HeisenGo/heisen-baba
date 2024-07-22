package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	ServerAddress          string `mapstructure:"SERVER_ADDRESS"`
	DSN                    string `mapstructure:"DSN"`
	TokenExpMinutes        uint   `mapstructure:"TOKEN_EXP_MINUTES"`
	RefreshTokenExpMinutes uint   `mapstructure:"REFRESH_TOKEN_EXP_MINUTES"`
	Secret                 string `mapstructure:"SECRET"`
}

func LoadConfig() (Config, error) {
	var cfg Config
	viper.AddConfigPath(".")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return cfg, err
	}

	err = viper.Unmarshal(&cfg)
	return cfg, err
}
