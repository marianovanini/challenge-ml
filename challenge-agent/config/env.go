package config

import (
	"log"

	"github.com/spf13/viper"
)

// Config contains application wide configurations

// Load recovers the env file and fils a config object
func Load(path string) {
	viper.SetConfigName("app")
	viper.AddConfigPath(path)
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err, "An error ocurring when loading configuration, using env variables instead")
	}

	log.Print("Configuration file loaded.")
}
