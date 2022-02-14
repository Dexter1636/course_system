package main

import (
	"context"
	"course_system/common"
	"github.com/spf13/viper"
)

func main() {
	common.InitConfig("test")
	common.InitDb()
	common.InitRdb(context.Background())
	r := RegisterRouter()
	port := viper.GetString("server.port")
	r.Run("0.0.0.0:" + port) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
