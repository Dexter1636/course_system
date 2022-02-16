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
)

type ICourseRedisRepository interface {
	CreateCourse(course *model.Course) (code vo.ErrNo)
	GetCourseById(id int64, course *model.Course) (code vo.ErrNo)
	GetAvailByCourseId(id int64, avail *int) (code vo.ErrNo)
	GetCourseListByStudentId(stuId int64, courseList *[]model.Course) (code vo.ErrNo)
}

type CourseRedisRepository struct {
	RDB *redis.Client
	ctx context.Context
}

func NewCourseRedisRepository() ICourseRedisRepository {
	return CourseRedisRepository{RDB: common.GetRDB(), ctx: common.GetCtx()}
}

func (crr CourseRedisRepository) CreateCourse(course *model.Course) (code vo.ErrNo) {
	code = vo.OK

	courseJson, err := json.Marshal(course)
	if err != nil {
		log.Println("Marshal ERROR when CreateCourse")
		code = vo.UnknownError
		return
	}
	_, err = crr.RDB.Set(crr.ctx, fmt.Sprintf("course:%d", course.Id), courseJson, 0).Result()
	if err != nil {
		log.Println("Redis set ERROR when CreateCourse")
		code = vo.UnknownError
		return
	}

	return code
}

func (crr CourseRedisRepository) GetCourseById(id int64, course *model.Course) (code vo.ErrNo) {
	code = vo.OK

	val, err := crr.RDB.Get(crr.ctx, fmt.Sprintf("course:%d", id)).Result()
	if err == redis.Nil {
		code = vo.CourseNotExisted
		return
	} else if err != nil {
		log.Println("Redis ERROR when GetCourseById")
		code = vo.UnknownError
		return
	}

	if err := json.Unmarshal([]byte(val), course); err != nil {
		log.Println("Unmarshal ERROR when GetCourseById")
		code = vo.UnknownError
		return
	}

	return code
}

func (crr CourseRedisRepository) GetAvailByCourseId(id int64, avail *int) (code vo.ErrNo) {
	code = vo.OK

	course := model.Course{}
	code = crr.GetCourseById(id, &course)

	if code == vo.OK {
		*avail = course.Avail
	}

	return code
}

func (crr CourseRedisRepository) GetCourseListByStudentId(stuId int64, courseList *[]model.Course) (code vo.ErrNo) {

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
		local enabled = user["Enabled"]
		if role ~= '2' or enabled ~= 1 then
			return {11}
		end

		local sc_exist = redis.call("EXISTS", sc_key)
		if sc_exist == 0 then
			return {13}
		end

		local sc = redis.call("SMEMBERS", sc_key)
		if #sc == 0 then
			return {13}
		end

		local course_list = {}
		for i = 1, #sc, 1 do
			local course = redis.call("GET", "course:" .. tostring(sc[i]))
    		table.insert(course_list, course)
		end

		return {0, cjson.encode(course_list)}
		`)
	val, err := lua.Run(crr.ctx, crr.RDB, keys, nil).Slice()
	if err != nil {
		log.Println(err.Error())
	}
	log.Println(val)

	// convert val[0] to ErrNo
	log.Println(val[0])
	resCode, ok := val[0].(int64)
	if !ok {
		log.Printf("Error when converting lua result to ErrNo in GetCourseListByStudentId")
		code = vo.UnknownError
		return
	}
	code = vo.ErrNo(resCode)

	// convert val[1] to courseList
	if code == vo.OK {
		courseListJson, ok := val[1].(string)
		if !ok {
			log.Printf("Error when converting lua result value to string in GetCourseListByStudentId")
			code = vo.UnknownError
			return
		}
		courseJsonList := make([]string, 8)
		err := json.Unmarshal([]byte(courseListJson), &courseJsonList)
		if err != nil {
			log.Println(err.Error())
			code = vo.UnknownError
			return
		}
		for _, courseJson := range courseJsonList {
			course := model.Course{}
			err := json.Unmarshal([]byte(courseJson), &course)
			if err != nil {
				log.Println(err.Error())
				code = vo.UnknownError
				return
			}
			*courseList = append(*courseList, course)
		}
	}

	return code
}
