package config

import "github.com/spf13/viper"

type Config struct {
	DB *DatabaseConfig `json:"db"`
}

type DatabaseConfig struct {
	DSN          string `json:"-"`
	MaxOpenConns int    `json:"max-open-conns"`
	MaxIdleConns int    `json:"max-idle-conns"`
	MaxIdleTime  string `json:"max-idle-time"`
	Timeout      int    `json:"timeout"`
}

func New() (*Config, error) {
	viper.AutomaticEnv()
	viper.AllowEmptyEnv(false)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME/.config/bookshelf")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	var config Config
	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
