package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Redis Redis
}

type Redis struct {
	Addr     string
	Password string
	DB       int
}

func LoadConfig() *Config {
	v := viper.New()

	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("./config/")

	if err := v.ReadInConfig(); err != nil {
		log.Fatal(err)
	}

	var config Config
	if err := v.Unmarshal(&config); err != nil {
		log.Fatal(err)
	}

	return &config
}
