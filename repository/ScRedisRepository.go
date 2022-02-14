package repository

import (
	"context"
	"course_system/common"
	"course_system/model"
	"course_system/vo"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
)

type IScRedisRepository interface {
	DeleteSc(stuId int64, courseId int64) (code vo.ErrNo)
	GetCourseListByStudentId(stuId int64, courseList *[]model.Course) (code vo.ErrNo)
}

type ScRedisRepository struct {
	RDB *redis.Client
	ctx context.Context
}

func NewScRedisRepository() IScRedisRepository {
	return ScRedisRepository{RDB: common.GetRDB(), ctx: common.GetCtx()}
}

func (crr ScRedisRepository) DeleteSc(stuId int64, courseId int64) (code vo.ErrNo) {
	code = vo.OK

	_, err := crr.RDB.SRem(crr.ctx, fmt.Sprintf("sc:%d", stuId), courseId).Result()
	if err == redis.Nil {
		log.Println("Redis ERROR when DeleteSc")
		code = vo.UnknownError
	}

	return code
}

func (crr ScRedisRepository) GetCourseListByStudentId(stuId int64, courseList *[]model.Course) (code vo.ErrNo) {

	keys := []string{fmt.Sprintf("user:%d", stuId), fmt.Sprintf("sc:%d", stuId)}
	var lua = redis.NewScript(`
		local user_key = KEYS[1]
		local sc_key = KEYS[2]
		
		local user_row = redis.call("GET", user_key)
		if user_row == false then
			return {11}
		end
		local user = cjson.decode(user_row)
		local role = user["RoleId"]
		if role ~= 2 then
  			return {11}
		end

		local sc_exist = redis.call("EXISTS", sc_key)
		if sc_exist == 0 then
			return {13}
		end

		local sc = redis.call("SMEMBERS", sc_key)
		return {0, sc}
		`)
	val, err := lua.Run(crr.ctx, crr.RDB, keys, nil).Slice()
	if err != nil {
		log.Println(err.Error())
	}
	log.Println(val)

	return code
}
