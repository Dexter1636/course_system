package common

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"os"
)

func InitConfig() {
	workDir, _ := os.Getwd()
	viper.SetConfigName("application")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(workDir + "/config")
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf(err.Error()))
	}
	viper.SetDefault("environment", "dev")
	env := viper.GetString("environment")
	if env == "test" {
		viper.SetConfigName("test")
		if err := viper.ReadInConfig(); err != nil {
			panic(fmt.Errorf(err.Error()))
		}
		log.Println("main: use test config")
	}
}
