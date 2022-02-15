package repository

import (
	"context"
	"course_system/common"
	"course_system/model"
	"course_system/vo"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"log"
	"strconv"
)

type IUserRedisRepository interface {
	ValidateStudentByUuid(uuid int64) (code vo.ErrNo)
}

type UserRedisRepository struct {
	DB  *gorm.DB
	RDB *redis.Client
	ctx context.Context
}

func NewUserRedisRepository() IUserRedisRepository {
	return UserRedisRepository{DB: common.GetDB(), RDB: common.GetRDB(), ctx: common.GetCtx()}
}

func (srr UserRedisRepository) ValidateStudentByUuid(uuid int64) (code vo.ErrNo) {
	code = vo.OK

	val, err := srr.RDB.Get(srr.ctx, "user:"+strconv.FormatInt(uuid, 10)).Result()
	if err == redis.Nil {
		code = vo.StudentNotExisted
		return
	} else if err != nil {
		log.Println("Redis ERROR when ValidateStudentByUuid")
		code = vo.UnknownError
		return
	} else {
		var user model.User
		if err := json.Unmarshal([]byte(val), &user); err != nil {
			log.Println("Unmarshal ERROR when ValidateStudentByUuid")
			code = vo.UnknownError
			return
		}
		if user.RoleId != fmt.Sprintf("%d", vo.Student) || user.Enabled != 1 {
			code = vo.StudentNotExisted
			return
		}
	}
	return
}
