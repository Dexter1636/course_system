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

func (t CreateCourseTest) getReq() interface{} {
	return t.Req
}

func (t CreateCourseTest) getExpCode() interface{} {
	return t.ExpCode
}

func (t CreateCourseTest) getExpResp() interface{} {
	return t.ExpResp
}

func (t GetCourseTest) getReq() interface{} {
	return t.Req
}

func (t GetCourseTest) getExpCode() interface{} {
	return t.ExpCode
}

func (t GetCourseTest) getExpResp() interface{} {
	return t.ExpResp
}

// ===============================

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

// ===============================
// ======== User ======== (Create && Get)

type CreateMemberTest struct {
	Req     vo.CreateMemberRequest
	ExpCode int
	ExpResp vo.CreateMemberResponse
}

type GetMemberTest struct {
	Req     vo.GetMemberRequest
	ExpCode int
	ExpResp vo.GetMemberResponse
}

func (t CreateMemberTest) getReq() interface{} {
	return t.Req
}

func (t CreateMemberTest) getExpCode() interface{} {
	return t.ExpCode
}

func (t CreateMemberTest) getExpResp() interface{} {
	return t.ExpResp
}

func (t GetMemberTest) getReq() interface{} {
	return t.Req
}

func (t GetMemberTest) getExpCode() interface{} {
	return t.ExpCode
}

func (t GetMemberTest) getExpResp() interface{} {
	return t.ExpResp
}
