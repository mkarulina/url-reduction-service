package config

import (
	"flag"
	"github.com/spf13/viper"
	"log"
)

func GetConfig(variable string) string {
	var foundFlag string

	flag.Visit(func(f *flag.Flag) {
		if Flags[f.Name] == variable {
			foundFlag = f.Value.String()
		}
	})
	if foundFlag == "" {
		viper.AddConfigPath("./config/")
		viper.SetConfigName("config")
		if err := viper.ReadInConfig(); err != nil {
			log.Println(err)
		}
		foundFlag = viper.GetString(variable)
	}
	return foundFlag
}
