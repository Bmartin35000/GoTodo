package config

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func LoadConfig() (config Config) { //TODO handle docker launch --> see viper.AutomaticEnv()
	viper.AddConfigPath("config")
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
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
