package controller

import (
	"course_system/dto"
	"course_system/model"
	"course_system/repository"
	"course_system/vo"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

type ICourseCommonController interface {
	CreateCourse(c *gin.Context)
	GetCourse(c *gin.Context)
}

type CourseCommonController struct {
	repo repository.ICourseRepository
}

func NewCourseCommonController() ICourseCommonController {
	return CourseCommonController{repo: repository.NewCourseRepository()}
}

func (ctl CourseCommonController) CreateCourse(c *gin.Context) {
	var req vo.CreateCourseRequest

	// validate data
	if err := c.ShouldBindJSON(&req); err != nil {
		panic(err.Error())
	}

	// create course
	course := model.Course{
		Name:  req.Name,
		Cap:   req.Cap,
		Avail: req.Cap,
	}

	code := ctl.repo.CreateCourse(&course)

	// response
	c.JSON(http.StatusOK, vo.CreateCourseResponse{
		Code: code,
		Data: struct {
			CourseID string
		}{CourseID: strconv.FormatInt(course.Id, 10)},
	})

}

func (ctl CourseCommonController) GetCourse(c *gin.Context) {
	var req vo.GetCourseRequest
	var code vo.ErrNo
	var course model.Course

	// response
	defer func() {
		c.JSON(http.StatusOK, vo.GetCourseResponse{
			Code: code,
			Data: dto.ToTCourse(course),
		})
	}()

	// validate data
	if err := c.ShouldBindJSON(&req); err != nil {
		panic(err.Error())
	}
	courseId, err := strconv.ParseInt(req.CourseID, 10, 64)
	if err != nil {
		code = vo.ParamInvalid
	}

	log.Println(req)

	// get course
	code = ctl.repo.GetCourseById(courseId, &course)
}
