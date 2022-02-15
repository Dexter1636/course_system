package common

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"io"
	"log"
	"os"
	"path"
	"time"
)

var env string

func InitLogger() {
	env = viper.GetString("environment")

	timeStr := time.Now().Format("20060102-150405")

	if env == "production" {
		wd, err := os.Getwd()
		if err != nil {
			panic(err.Error())
		}
		logFilePath := wd + "/logs/"
		logFileName := timeStr + ".log"
		fileName := path.Join(logFilePath, logFileName)

		if err := os.MkdirAll(logFilePath, 0777); err != nil {
			panic(err.Error())
		}

		gin.DisableConsoleColor()
		f, err := os.Create(fileName)
		if err != nil {
			panic(err.Error())
		}
		gin.DefaultWriter = io.MultiWriter(f)
		log.SetOutput(f)

		fmt.Println("==== logger env: production ====")
	}
}
