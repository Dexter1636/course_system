package common

import (
	"context"
	"course_system/model"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"time"
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
	DB = db
	fmt.Println("Connected to database.")
}

func MessageQueue() {
	for {
		count, err := RDB.LLen(Ctx, "MessageQueue").Result()
		if err != nil {
			log.Println(err.Error())
		}
		if count > 0 {
			val, err := RDB.LPop(Ctx, "MessageQueue").Result()
			if err != nil {
				log.Println(err.Error())
			}
			var sc model.Sc
			if err = json.Unmarshal([]byte(val), &sc); err != nil {
				log.Println(err.Error())
			}
			err = DB.Transaction(func(tx *gorm.DB) error {
				// check avail
				course := model.Course{Id: sc.CourseId}
				if err := tx.Select("avail").Take(&course, sc.CourseId).Error; err != nil {
					return err
				}
				// update avail
				course.Avail--
				if err := tx.Model(&course).Update("avail", course.Avail).Error; err != nil {
					return err
				}
				// create sc record
				if err := tx.Create(&sc).Error; err != nil {
					return err
				}
				return nil
			})
			if err != nil {
				log.Println(err.Error())
			}
		} else {
			time.Sleep(1 * time.Second)
			fmt.Println("Sleep 1 second")
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
	//初始化消息队列
	go MessageQueue()
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
