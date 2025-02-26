package config

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func LoadConfig() (config Config) {
	viper.AddConfigPath("config")
	viper.SetConfigType("yaml")
	if os.Getenv("source") == "docker" {
		viper.SetConfigName("config.docker")
	} else {
		viper.SetConfigName("config")
	}
	err := viper.ReadInConfig()
	if err != nil {
		log.WithFields(log.Fields{"details": err}).Panic("failed to load env file")
		panic("failed to load env file")
	}
	err = viper.Unmarshal(&config)
	if err != nil {
		log.WithFields(log.Fields{"details": err}).Panic("failed to unparse env file")
		panic("failed to unparse env file")
	}
	return // automatic return by name
}
