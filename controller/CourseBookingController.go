package controller

import (
	"context"
	"course_system/common"
	"course_system/model"
	"course_system/repository"
	"course_system/vo"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"log"
	"net/http"
	"strconv"
)

type ICourseBookingController interface {
	BookCourse(c *gin.Context)
	GetStudentCourse(c *gin.Context)
}

type CourseBookingController struct {
	repo            repository.ICourseRepository
	courseRedisRepo repository.ICourseRedisRepository
	scRedisRepo     repository.IScRedisRepository
	userRedisRepo   repository.IUserRedisRepository
	DB              *gorm.DB
	RDB             *redis.Client
	ctx             context.Context
}

func NewCourseBookingController() ICourseBookingController {
	return CourseBookingController{
		repo:            repository.NewCourseRepository(),
		courseRedisRepo: repository.NewCourseRedisRepository(),
		scRedisRepo:     repository.NewScRedisRepository(),
		userRedisRepo:   repository.NewUserRedisRepository(),
		DB:              common.GetDB(),
		RDB:             common.GetRDB(),
		ctx:             common.GetCtx(),
	}
}

func (ctl CourseBookingController) BookCourse(c *gin.Context) {
	var req vo.BookCourseRequest
	code := vo.OK

	// response
	defer func() { c.JSON(http.StatusOK, vo.BookCourseResponse{Code: code}) }()

	// validate data
	if err := c.ShouldBindJSON(&req); err != nil {
		code = vo.ParamInvalid
		return
	}
	studentId, err := strconv.ParseInt(req.StudentID, 10, 64)
	if err != nil {
		code = vo.ParamInvalid
		return
	}
	courseId, err := strconv.ParseInt(req.CourseID, 10, 64)
	if err != nil {
		code = vo.ParamInvalid
		return
	}

	// =============================================================
	// book course (v2: cache "course", "sc" and "student" to Redis)
	// 1. validate student
	// 2. validate course
	// 3. validate sc
	// 4. update course avail - 1
	// 5. write new data to MySQL
	//    - if failed:
	//   	1. delete sc
	//  	2. update course avail + 1
	// =============================================================

	// 1. validate student
	code = ctl.userRedisRepo.ValidateStudentByUuid(studentId)
	if code != vo.OK {
		log.Println("[BookCourse] ErrNo: ", code)
		return
	}
	// ==== do step 2, 3 and 4 in Lua script ====
	// 2. validate course
	// 3. validate sc
	// 4. update course avail - 1
	keys := []string{fmt.Sprintf("course:%d", courseId), fmt.Sprintf("sc:%d", studentId)}
	values := []interface{}{courseId}
	var lua = redis.NewScript(`
		local course_key = KEYS[1]
		local sc_key = KEYS[2]
		local course_id = ARGV[1]
		
		local course_row = redis.call("GET", course_key)
		if course_row == false then
			return 12
		end
		local course = cjson.decode(course_row)
		local avail = course["Avail"]
		if avail <= 0 then
  			return 7
		end

		local sc_value_exist = redis.call("SISMEMBER", sc_key, course_id)
		if sc_value_exist ~= 0 then
			return 14
		end

		course["Avail"] = avail - 1
		course_row = cjson.encode(course)
		redis.call("SET", course_key, course_row)
		redis.call("SADD", sc_key, course_id)
		
		return 0
		`)
	codeInt, err := lua.Run(ctl.ctx, ctl.RDB, keys, values...).Int()
	if err != nil {
		log.Println(err.Error())
	}
	fmt.Println(codeInt)
	code = vo.ErrNo(codeInt)

	// 5. write new data to MySQL
	if code == vo.OK {
		err = ctl.DB.Transaction(func(tx *gorm.DB) error {
			// check avail
			course := model.Course{Id: courseId}
			if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Select("avail").First(&course, courseId).Error; err != nil {
				log.Println(err.Error())
				if errors.Is(err, gorm.ErrRecordNotFound) {
					code = vo.CourseNotExisted
				} else {
					code = vo.UnknownError
				}
				return err
			}
			if course.Avail <= 0 {
				code = vo.CourseNotAvailable
				return errors.New("CourseNotAvailable")
			}
			// update avail
			course.Avail--
			if err := tx.Model(&course).Update("avail", course.Avail).Error; err != nil {
				log.Println(err.Error())
				code = vo.UnknownError
				return err
			}
			// create sc record
			sc := model.Sc{
				StudentId: studentId,
				CourseId:  courseId,
			}
			if err := tx.Create(&sc).Error; err != nil {
				log.Println(err.Error())
				var mysqlErr *mysql.MySQLError
				if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 { // student already have this course
					code = vo.StudentHasCourse
				} else {
					code = vo.UnknownError
				}
				return err
			}
			return nil
		})
		if err != nil {
			log.Println(err.Error())
			// rollback Redis data
			if resCode := ctl.scRedisRepo.DeleteSc(studentId, courseId); resCode == vo.UnknownError {
				code = resCode
			}
			return
		}
	}

}

func (ctl CourseBookingController) GetStudentCourse(c *gin.Context) {
	var req vo.GetStudentCourseRequest
	code := vo.OK
	courseList := make([]model.Course, 0, 8)
	tCourseList := make([]vo.TCourse, 0, 8)

	// response
	defer func() {
		c.JSON(http.StatusOK, vo.GetStudentCourseResponse{
			Code: code,
			Data: struct {
				CourseList []vo.TCourse
			}{tCourseList},
		})
		log.Printf("code: %d\n\n", code)
	}()

	// validate data
	if err := c.ShouldBindQuery(&req); err != nil {
		code = vo.ParamInvalid
		return
	}
	studentId, err := strconv.ParseInt(req.StudentID, 10, 64)
	if err != nil {
		code = vo.ParamInvalid
		return
	}

	// get course
	code = ctl.courseRedisRepo.GetCourseListByStudentId(studentId, &courseList)

	// convert query result to response type
	for _, course := range courseList {
		tCourseList = append(tCourseList, vo.TCourse{
			CourseID:  strconv.FormatInt(course.Id, 10),
			Name:      course.Name,
			TeacherID: strconv.FormatInt(course.TeacherId, 10),
		})
	}
}
