package main

import (
	"auth/command"
	"auth/common/logger"
	"auth/config"
	"fmt"
	"os"

	"github.com/spf13/viper"
)

func readConfig() {
	var err error

	viper.SetConfigFile("base.env")
	viper.SetConfigType("props")
	err = viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
		return
	}

	if _, err := os.Stat("base.env"); os.IsNotExist(err) {
		fmt.Println("WARNING: file base.env not found")
	} else {
		viper.SetConfigFile("base.env")
		viper.SetConfigType("props")
		err = viper.MergeInConfig()
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	err = viper.Unmarshal(&config.Config)
	for v, key := range viper.AllKeys() {
		fmt.Println(key, " ", v)
		viper.BindEnv(key)
	}

	for v, key := range viper.AllKeys() {
		fmt.Println(key, " ", v)
		viper.BindEnv(key)
	}
}

func main() {
	readConfig()
	raventClient := logger.NewRavenClient()
	logger.NewLogger(raventClient)
	command.Execute()
}
