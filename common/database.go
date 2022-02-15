package common

import (
	"context"
	"course_system/model"
	"encoding/json"
	"errors"
	"fmt"
  "github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB
var RDB *redis.Client
var Ctx context.Context

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
	// set connection pool size
	sqlDB, err := db.DB()
	if err != nil {
		panic("failed to config db connection pool, err: " + err.Error())
	}
	sqlDB.SetMaxOpenConns(150)
	DB = db
	fmt.Println("Connected to database.")
	
	// @Author 彭守恒 2022-02-15 02:45 删除以避免循环依赖
	// data.CheckAdmin()
	var u model.User
	if err := DB.Where("user_name = ?", "JudgeAdmin").Take(&u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			GetDB().Exec("INSERT INTO user(uuid, user_name, nick_name, password, role_id, enabled) " +
				"VALUES (1, 'JudgeAdmin', 'JudgeAdmin', 'JudgePassword2022', '1', 1)")
		}
	}
}

func InitRedisData() {
	//清空数据
	RDB.FlushDB(Ctx)
	//读取user表数据
	var users []model.User
	if err := DB.Find(&users).Error; err != nil {
		panic(err.Error())
	}
	for _, user := range users {
		val, err := json.Marshal(user)
		if err != nil {
			panic(err.Error())
		}
		err = RDB.Set(Ctx, fmt.Sprintf("user:%d", user.Uuid), val, 0).Err()
		if err != nil {
			panic(err.Error())
		}
	}
	//读取course表数据
	var courses []model.Course
	if err := DB.Find(&courses).Error; err != nil {
		panic(err.Error())
	}
	for _, course := range courses {
		val, err := json.Marshal(course)
		if err != nil {
			panic(err.Error())
		}
		err = RDB.Set(Ctx, fmt.Sprintf("course:%d", course.Id), val, 0).Err()
		if err != nil {
			panic(err.Error())
		}
	}
	//读取sc表数据
	var scs []model.Sc
	if err := DB.Find(&scs).Error; err != nil {
		panic(err.Error())
	}
	for _, sc := range scs {
		err := RDB.SAdd(Ctx, fmt.Sprintf("sc:%d", sc.StudentId), sc.CourseId, 0).Err()
		if err != nil {
			panic(err.Error())
		}
	}
}

func InitRdb(ctx context.Context) {
	host := viper.GetString("redis.host")
	port := viper.GetString("redis.port")
	db := viper.GetString("redis.db")
	user := viper.GetString("redis.user")
	password := viper.GetString("redis.password")
	redisUrl := fmt.Sprintf("redis://%s:%s@%s:%s/%s", user, password, host, port, db)
	opt, err := redis.ParseURL(redisUrl)
	if err != nil {
		panic(err)
	}
	rdb := redis.NewClient(opt)
	if pong, err := rdb.Ping(ctx).Result(); err != nil || pong != "PONG" {
		panic("failed to connect to redis server, err: " + err.Error())
	}
	RDB = rdb
	Ctx = ctx
	fmt.Println("Connected to redis server.")
	InitRedisData()
	fmt.Println("redis data initialization is complete.")
}

func GetDB() *gorm.DB {
	return DB
}

func GetRDB() *redis.Client {
	return RDB
}

func GetCtx() context.Context {
	return Ctx
}