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
	repo            repository.ICourseRepository
	courseRedisRepo repository.ICourseRedisRepository
}

func NewCourseCommonController() ICourseCommonController {
	return CourseCommonController{
		repo:            repository.NewCourseRepository(),
		courseRedisRepo: repository.NewCourseRedisRepository(),
	}
}

func (ctl CourseCommonController) CreateCourse(c *gin.Context) {
	var req vo.CreateCourseRequest
	var course model.Course
	var code vo.ErrNo

	// response
	defer func() {
		c.JSON(http.StatusOK, vo.CreateCourseResponse{
			Code: code,
			Data: struct {
				CourseID string
			}{CourseID: strconv.FormatInt(course.Id, 10)},
		})
		log.Printf("[CreateCourse] code: %d\n", code)
	}()

	// validate data
	if err := c.ShouldBindJSON(&req); err != nil {
		code = vo.ParamInvalid
		return
	}

	// course instance
	course = model.Course{
		Name:  req.Name,
		Cap:   req.Cap,
		Avail: req.Cap,
	}

	// create course in MySQL
	code = ctl.repo.CreateCourse(&course)

	// create course in Redis
	code = ctl.courseRedisRepo.CreateCourse(&course)
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
		log.Printf("[GetCourse] code: %d\n", code)
	}()

	// validate data
	if err := c.ShouldBindQuery(&req); err != nil {
		// TODO
		panic(err.Error())
	}
	courseId, err := strconv.ParseInt(req.CourseID, 10, 64)
	if err != nil {
		code = vo.ParamInvalid
	}

	// get course
	code = ctl.courseRedisRepo.GetCourseById(courseId, &course)
}
