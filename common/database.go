package common

import (
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDb() {
	// Capture connection properties
	//driverName := viper.GetString("datasource.driverName")
	host := viper.GetString("datasource.host")
	port := viper.GetString("datasource.port")
	database := viper.GetString("datasource.database")
	username := viper.GetString("datasource.username")
	password := viper.GetString("datasource.password")
	charset := viper.GetString("datasource.charset")
	loggerLevel := viper.GetString("logger.level")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true",
		username, password, host, port, database, charset)
	// config
	config := &gorm.Config{}
	if loggerLevel == "info" {
		config.Logger = logger.Default.LogMode(logger.Info)
	}
	// Get a database handle.
	db, err := gorm.Open(mysql.Open(dsn), config)
	if err != nil {
		panic("failed to connect to database, err: " + err.Error())
	}
	DB = db
	fmt.Println("Connected to database.")
}

func GetDB() *gorm.DB {
	return DB
}
