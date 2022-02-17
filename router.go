package main

import (
	"course_system/controller"
	"course_system/vo"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func RegisterRouter() *gin.Engine {
	r := gin.Default()

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	r.Use(gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			c.JSON(http.StatusInternalServerError, vo.ResponseMeta{Code: vo.UnknownError})
			log.Println("========================================")
			log.Printf("!!PANIC!! ERR: %s\n", err)
			log.Println("========================================")
		}
		c.AbortWithStatus(http.StatusInternalServerError)
	}))

	g := r.Group("/api/v1")

	// ping test
	g.GET("/ping", controller.Ping)

	// 成员管理
	uc := controller.NewUserController()
	g.POST("/member/create", uc.Create)
	g.GET("/member", uc.Member)
	g.GET("/member/list", uc.List)
	g.POST("/member/update", uc.Update)
	g.POST("/member/delete", uc.Delete)

	// 登录
	ac := controller.NewAuthController()
	g.POST("/auth/login", ac.Login)
	g.POST("/auth/logout", ac.Logout)
	g.GET("/auth/whoami", ac.WhoAmI)

	// 排课
	ccc := controller.NewCourseCommonController()
	g.POST("/course/create", ccc.CreateCourse)
	g.GET("/course/get", ccc.GetCourse)
	csc := controller.NewCourseScheduleController()
	g.POST("/teacher/bind_course", csc.Bind)
	g.POST("/teacher/unbind_course", csc.Unbind)
	g.GET("/teacher/get_course", csc.Get)
	g.POST("/course/schedule", csc.Schedule)

	// 抢课
	cbc := controller.NewCourseBookingController()
	g.POST("/student/book_course", cbc.BookCourse)
	g.GET("/student/course", cbc.GetStudentCourse)

	return r
}
