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
	ctx context.Context
}

func NewScRedisRepository() IScRedisRepository {
	return ScRedisRepository{RDB: common.GetRDB(), ctx: common.GetCtx()}
}

func (srr ScRedisRepository) DeleteSc(stuId int64, courseId int64) (code vo.ErrNo) {
	code = vo.OK

	_, err := srr.RDB.SRem(srr.ctx, fmt.Sprintf("sc:%d", stuId), courseId).Result()
	if err == redis.Nil {
		log.Println("Redis ERROR when DeleteSc")
		code = vo.UnknownError
	}

	return code
}
