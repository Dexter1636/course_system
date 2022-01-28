package main

import (
	"course_system/controller"
	"github.com/gin-gonic/gin"
)

func RegisterRouter(r *gin.Engine) {
	g := r.Group("/api/v1")

	// ping test
	g.GET("/ping", controller.Ping)

	// 成员管理
	g.POST("/member/create")
	g.GET("/member")
	g.GET("/member/list")
	g.POST("/member/update")
	g.POST("/member/delete")

	// 登录
	ac := controller.NewAuthController()
	g.POST("/auth/login", ac.Login)
	g.POST("/auth/logout", ac.Logout)
	g.GET("/auth/whoami", ac.WhoAmI)

	// 排课
	ccc := controller.NewCourseCommonController()
	g.POST("/course/create", ccc.CreateCourse)
	g.GET("/course/get", ccc.GetCourse)

	g.POST("/teacher/bind_course")
	g.POST("/teacher/unbind_course")
	g.GET("/teacher/get_course")
	g.POST("/course/schedule")

	// 抢课
	g.POST("/student/book_course")
	g.GET("/student/course")

}
