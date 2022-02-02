package controller

import (
	"course_system/common"
	"course_system/model"
	"course_system/vo"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

type ICourseBookingController interface {
	Bind(c *gin.Context)
	Unbind(c *gin.Context)
}
type CourseBookingController struct {
	DB *gorm.DB
}

func NewCourseBookingController() ICourseBookingController {
	db := common.GetDB()
	return CourseBookingController{DB: db}
}
func (ctl CourseBookingController) Bind(c *gin.Context) {
	var req vo.BindCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		panic(err.Error())
	}
	var sample model.Course
	a := ctl.DB.Model(&model.Course{}).First(&sample, req.CourseID)
	if a.Error == gorm.ErrEmptySlice {
		c.JSON(http.StatusOK, vo.BindCourseResponse{Code: vo.CourseNotExisted})
	} else if sample.TeacherId != 0 {
		c.JSON(http.StatusOK, vo.BindCourseResponse{Code: vo.CourseHasBound})
	} else {
		ctl.DB.Model(&model.Course{}).First(sample, req.CourseID).Update("TeacherId", req.TeacherID)
		c.JSON(http.StatusOK, vo.BindCourseResponse{Code: vo.OK})
	}
}
func (ctl CourseBookingController) Unbind(c *gin.Context) {
	var req vo.UnbindCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		panic(err.Error())
	}
	var sample model.Course
	a := ctl.DB.Model(&model.Course{}).First(&sample, req.CourseID)
	if a.Error == gorm.ErrEmptySlice {
		c.JSON(http.StatusOK, vo.UnbindCourseResponse{Code: vo.CourseNotExisted})
	} else if sample.TeacherId == 0 {
		c.JSON(http.StatusOK, vo.UnbindCourseResponse{Code: vo.CourseNotBind})
	} else {
		ctl.DB.Model(&model.Course{}).First(sample, req.CourseID).Update("TeacherId", 0)
		c.JSON(http.StatusOK, vo.UnbindCourseResponse{Code: vo.OK})
	}
}
