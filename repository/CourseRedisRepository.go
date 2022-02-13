package repository

import (
	"context"
	"course_system/common"
	"course_system/vo"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"strconv"
)

type ICourseRedisRepository interface {
	GetAvailByCourseId(id int64, avail *int) (code vo.ErrNo)
}

type CourseRedisRepository struct {
	RDB *redis.Client
	Ctx context.Context
}

func NewCourseRedisRepository() ICourseRedisRepository {
	return CourseRedisRepository{RDB: common.GetRDB(), Ctx: common.GetCtx()}
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
