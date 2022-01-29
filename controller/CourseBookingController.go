package controller

import (
	"course_system/common"
	"course_system/model"
	"course_system/vo"
	"errors"
	"github.com/gin-gonic/gin"
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
	DB *gorm.DB
}

func NewCourseBookingController() ICourseBookingController {
	db := common.GetDB()
	return CourseBookingController{DB: db}
}

func (ctl CourseBookingController) BookCourse(c *gin.Context) {
	var req vo.BookCourseRequest
	code := vo.OK

	// response
	defer c.JSON(http.StatusOK, vo.BookCourseResponse{Code: code})

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

	// book course (v1: select for update)
	// 1. check avail
	// 2. update avail
	// 3. create sc record
	err = ctl.DB.Transaction(func(tx *gorm.DB) error {
		// check avail
		course := model.Course{Id: courseId}
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Select("avail").First(&course, courseId).Error; err != nil {
			log.Println(err)
			return err
		}
		if course.Avail <= 0 {
			code = vo.CourseNotAvailable
			return errors.New("CourseNotAvailable")
		}
		// update avail
		if err := tx.Model(&course).Update("avail", course.Avail).Error; err != nil {
			log.Println(err)
			return err
		}
		// create sc record
		sc := model.Sc{
			StudentId: studentId,
			CourseId:  courseId,
		}
		if err := tx.Create(&sc).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return
	}
}

func (ctl CourseBookingController) GetStudentCourse(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}
