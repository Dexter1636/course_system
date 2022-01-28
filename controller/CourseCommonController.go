package controller

import (
	"course_system/common"
	"course_system/model"
	"course_system/vo"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
)

type ICourseCommonController interface {
	CreateCourse(c *gin.Context)
	GetCourse(c *gin.Context)
}

type CourseCommonController struct {
	DB *gorm.DB
}

func NewCourseCommonController() ICourseCommonController {
	db := common.GetDB()
	return CourseCommonController{DB: db}
}

func (ctl CourseCommonController) CreateCourse(c *gin.Context) {
	var req vo.CreateCourseRequest

	// validate data
	if err := c.ShouldBindJSON(&req); err != nil {
		panic(err.Error())
	}

	// create course
	course := model.Course{
		Name: req.Name,
		Cap:  req.Cap,
	}

	if err := ctl.DB.Create(&course).Error; err != nil {
		panic(err.Error())
	}

	// response
	// TODO: [bug] wrong response type
	c.JSON(http.StatusOK, vo.GetCourseResponse{
		Code: 0,
		Data: vo.TCourse{
			CourseID:  strconv.FormatInt(course.Id, 10),
			Name:      course.Name,
			TeacherID: strconv.FormatInt(course.TeacherId, 10),
		},
	})

}

func (ctl CourseCommonController) GetCourse(c *gin.Context) {
	var req vo.GetCourseRequest

	// validate data
	if err := c.ShouldBindJSON(&req); err != nil {
		panic(err.Error())
	}

	log.Println(req)

	// get course
	var course model.Course
	if err := ctl.DB.First(&course, req.CourseID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusOK, vo.GetCourseResponse{Code: vo.CourseNotExisted})
			return
		} else {
			panic(err.Error())
		}
	}

	// response
	c.JSON(http.StatusOK, vo.GetCourseResponse{
		Code: vo.OK,
		Data: vo.TCourse{
			CourseID:  strconv.FormatInt(course.Id, 10),
			Name:      course.Name,
			TeacherID: strconv.FormatInt(course.TeacherId, 10),
		},
	})

}
