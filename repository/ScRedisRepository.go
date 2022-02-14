package repository

import (
	"context"
	"course_system/common"
	"course_system/vo"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
)

type IScRedisRepository interface {
	DeleteSc(stuId int64, courseId int64) (code vo.ErrNo)
}

type ScRedisRepository struct {
	RDB *redis.Client
	Ctx context.Context
}

func NewScRedisRepository() IScRedisRepository {
	return ScRedisRepository{RDB: common.GetRDB(), Ctx: common.GetCtx()}
}

func (crr ScRedisRepository) DeleteSc(stuId int64, courseId int64) (code vo.ErrNo) {
	code = vo.OK

	_, err := crr.RDB.SRem(crr.Ctx, fmt.Sprintf("sc:%d", stuId), courseId).Result()
	if err == redis.Nil {
		log.Println("Redis ERROR when DeleteSc")
		code = vo.UnknownError
	}

	return code
}
