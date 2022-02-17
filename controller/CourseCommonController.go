package controller

import (
	"course_system/dto"
	"course_system/model"
	"course_system/repository"
	"course_system/utils"
	"course_system/vo"
	"github.com/gin-gonic/gin"
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
	var resp vo.CreateCourseResponse
	var course model.Course
	var code vo.ErrNo

	// response
	defer func() {
		resp = vo.CreateCourseResponse{
			Code: code,
			Data: struct {
				CourseID string
			}{CourseID: strconv.FormatInt(course.Id, 10)},
		}
		c.JSON(http.StatusOK, resp)
		utils.LogReqRespBody(req, resp, "CreateCourse")
	}()

	// validate data
	if err := c.ShouldBindJSON(&req); err != nil {
		code = vo.ParamInvalid
		return
	}
	if code = vo.ValidateCreateCourseReq(req); code != vo.OK {
		return
	}

	// log request body
	utils.LogBody(req, "CreateCourse.req")

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
	var resp vo.GetCourseResponse
	var code vo.ErrNo
	var course model.Course

	// response
	defer func() {
		resp = vo.GetCourseResponse{
			Code: code,
			Data: dto.ToTCourse(course),
		}
		c.JSON(http.StatusOK, resp)
		utils.LogReqRespBody(req, resp, "GetCourse")
	}()

	// validate data
	if err := c.ShouldBindQuery(&req); err != nil {
		code = vo.ParamInvalid
		return
	}
	if code = vo.ValidateGetCourseReq(req); code != vo.OK {
		return
	}

	courseId, _ := strconv.ParseInt(req.CourseID, 10, 64)

	// get course
	code = ctl.courseRedisRepo.GetCourseById(courseId, &course)
}
