package controller

import (
	"course_system/common"
	"course_system/vo"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ICourseBookingController interface {
	BookCourse(c *gin.Context)
	GetCourseList(c *gin.Context)
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

	// book course

}

func (ctl CourseBookingController) GetCourseList(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}
