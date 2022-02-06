package test

import (
	"course_system/vo"
)

type BaseTest interface {
	getReq() interface{}
	getExpCode() interface{}
	getExpResp() interface{}
}

// ======== CourseCommon ========

type CreateCourseTest struct {
	Req     vo.CreateCourseRequest
	ExpCode int
	ExpResp vo.CreateCourseResponse
}

type GetCourseTest struct {
	Req     vo.GetCourseRequest
	ExpCode int
	ExpResp vo.GetCourseResponse
}

// ======== CourseBooking ========

type BookCourseTest struct {
	Req     vo.BookCourseRequest
	ExpCode int
	ExpResp vo.BookCourseResponse
}

type GetStudentCourseTest struct {
	Req     vo.GetStudentCourseRequest
	ExpCode int
	ExpResp vo.GetStudentCourseResponse
}

func (g GetStudentCourseTest) getReq() interface{} {
	return g.Req
}

func (g GetStudentCourseTest) getExpCode() interface{} {
	return g.ExpCode
}

func (g GetStudentCourseTest) getExpResp() interface{} {
	return g.ExpResp
}
