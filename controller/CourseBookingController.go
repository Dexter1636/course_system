package controller

import (
	"course_system/common"
	"course_system/model"
	"course_system/repository"
	"course_system/vo"
	"errors"
	"github.com/gin-gonic/gin"
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
	repo repository.ICourseRepository
	DB   *gorm.DB
}

func NewCourseBookingController() ICourseBookingController {
	return CourseBookingController{
		repo: repository.NewCourseRepository(),
		DB:   common.GetDB(),
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

	// book course (v1.1: select for update)
	// 1. validate student
	// 2. check course avail
	// 3. update course avail
	// 4. create sc record
	err = ctl.DB.Transaction(func(tx *gorm.DB) error {
		// validate student
		var count int64
		if err := tx.Model(&model.User{}).Where("uuid = ?", studentId).Count(&count).Error; err != nil {
			log.Println(err.Error())
			code = vo.UnknownError
			return err
		}
		if count <= 0 {
			code = vo.StudentNotExisted
			return errors.New("StudentNotExisted")
		}
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
		return
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
	}()

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

	// get course
	code = ctl.repo.GetCourseListByStudentId(studentId, &courseList)

	// convert query result to response type
	for _, course := range courseList {
		tCourseList = append(tCourseList, vo.TCourse{
			CourseID:  strconv.FormatInt(course.Id, 10),
			Name:      course.Name,
			TeacherID: strconv.FormatInt(course.TeacherId, 10),
		})
	}
}
