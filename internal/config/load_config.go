package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	DBParams struct {
		Address string `mapstructure:"db_addr"`
		Password string `mapstructure:"db_pass"`
		Number int `mapstructure:"db_num"`
		Protocol int `mapstructure:"db_protocol"`
	} `mapstructure:"db_params"`
	ServerParams struct {
		ListenPort string `mapstructure:"listen_port"`
	} `mapstructure:"server_params"`
}

func LoadConfig() *Config {
	viper.AddConfigPath("../internal/config")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}
	AppConfig := new(Config)
	viper.Unmarshal(AppConfig)
	return AppConfig
}
