package repository

import (
	"context"
	"course_system/common"
	"course_system/model"
	"course_system/vo"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"strconv"
)

type ICourseRedisRepository interface {
	CreateCourse(course model.Course) (code vo.ErrNo)
	GetAvailByCourseId(id int64, avail *int) (code vo.ErrNo)
}

type CourseRedisRepository struct {
	RDB *redis.Client
	Ctx context.Context
}

func NewCourseRedisRepository() ICourseRedisRepository {
	return CourseRedisRepository{RDB: common.GetRDB(), Ctx: common.GetCtx()}
}

func (crr CourseRedisRepository) CreateCourse(course model.Course) (code vo.ErrNo) {
	code = vo.OK

	courseJson, err := json.Marshal(course)
	if err != nil {
		log.Println("Unmarshal ERROR when CreateCourse")
		code = vo.UnknownError
		return
	}
	_, err = crr.RDB.Set(crr.Ctx, fmt.Sprintf("course:%d", course.Id), courseJson, 0).Result()
	if err != nil {
		log.Println("Redis set ERROR when CreateCourse")
		code = vo.UnknownError
		return
	}

	return code
}

func (crr CourseRedisRepository) GetAvailByCourseId(id int64, avail *int) (code vo.ErrNo) {
	code = vo.OK

	val, err := crr.RDB.Get(crr.Ctx, fmt.Sprintf("course:%d", id)).Result()
	if err == redis.Nil {
		code = vo.CourseNotExisted
	} else if err != nil {
		log.Println("Redis ERROR when GetAvailByCourseId")
		code = vo.UnknownError
	} else {
		availInt, _ := strconv.Atoi(val)
		*avail = availInt
	}

	return code
}
