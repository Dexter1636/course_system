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

func (t BookCourseTest) getReq() interface{} {
	return t.Req
}

func (t BookCourseTest) getExpCode() interface{} {
	return t.ExpCode
}

func (t BookCourseTest) getExpResp() interface{} {
	return t.ExpResp
}

func (t GetStudentCourseTest) getReq() interface{} {
	return t.Req
}

func (t GetStudentCourseTest) getExpCode() interface{} {
	return t.ExpCode
}

func (t GetStudentCourseTest) getExpResp() interface{} {
	return t.ExpResp
}
