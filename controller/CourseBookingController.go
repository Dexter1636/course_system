package controller

import (
	"course_system/common"
	"course_system/model"
	"course_system/vo"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"net/http"
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

	// validate data
	if err := c.ShouldBindJSON(&req); err != nil {
		panic(err.Error())
	}

	// book course (v1: select for update)
	// 1. check avail
	// 2. delete avail
	// 3. create sc record
	err := ctl.DB.Transaction(func(tx *gorm.DB) error {
		var course model.Course
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Select("avail").First(&course, req.CourseID); err != nil {
			return errors.New("")
		}
		//if course <= 0 {
		//	return errors.New("")
		//}
		if err := tx.Model(&course).Updates("avail"); err != nil {
			return errors.New("")
		}
		return nil
	})
	if err != nil {
		return
	}

	// response
	c.JSON(http.StatusOK, vo.BookCourseResponse{Code: vo.OK})
}

func (ctl CourseBookingController) GetStudentCourse(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}
