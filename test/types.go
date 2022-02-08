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

//<<<<<<< HEAD =======

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

// ======================================= CourseBooking ========

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

//=======================Schedule Course========================

type BindCourseTest struct {
	Req     vo.BindCourseRequest
	ExpCode int
	ExpResp vo.BindCourseResponse
}
type UnBindCourseTest struct {
	Req     vo.UnbindCourseRequest
	ExpCode int
	ExpResp vo.UnbindCourseResponse
}
type TGetCourseTests struct {
	Req     vo.GetTeacherCourseRequest
	ExpCode int
	ExpResp vo.GetTeacherCourseResponse
}
type ScheduleTest struct {
	Req     vo.ScheduleCourseRequest
	ExpCode int
	ExpResp vo.ScheduleCourseResponse
}

func (t BindCourseTest) getReq() interface{} {
	return t.Req
}

func (t BindCourseTest) getExpCode() interface{} {
	return t.ExpCode
}

func (t BindCourseTest) getExpResp() interface{} {
	return t.ExpResp
}

func (t UnBindCourseTest) getReq() interface{} {
	return t.Req
}

func (t UnBindCourseTest) getExpCode() interface{} {
	return t.ExpCode
}

func (t UnBindCourseTest) getExpResp() interface{} {
	return t.ExpResp
}

func (t TGetCourseTests) getReq() interface{} {
	return t.Req
}

func (t TGetCourseTests) getExpCode() interface{} {
	return t.ExpCode
}

func (t TGetCourseTests) getExpResp() interface{} {
	return t.ExpResp
}

func (t ScheduleTest) getReq() interface{} {
	return t.Req
}

func (t ScheduleTest) getExpCode() interface{} {
	return t.ExpCode
}

func (t ScheduleTest) getExpResp() interface{} {
	return t.ExpResp
}

// ======== User ========(other)

type GetMemberListTest struct {
	Req     vo.GetMemberListRequest
	ExpCode int
	ExpResp vo.GetMemberListResponse
}

type UpdateMemberTest struct {
	Req     vo.UpdateMemberRequest
	ExpCode int
	ExpResp vo.UpdateMemberResponse
}

type DeleteMemberTest struct {
	Req     vo.DeleteMemberRequest
	ExpCode int
	ExpResp vo.DeleteMemberResponse
}

func (t GetMemberListTest) getReq() interface{} {
	return t.Req
}

func (t GetMemberListTest) getExpCode() interface{} {
	return t.ExpCode
}

func (t GetMemberListTest) getExpResp() interface{} {
	return t.ExpResp
}

func (t UpdateMemberTest) getReq() interface{} {
	return t.Req
}

func (t UpdateMemberTest) getExpCode() interface{} {
	return t.ExpCode
}

func (t UpdateMemberTest) getExpResp() interface{} {
	return t.ExpResp
}

func (t DeleteMemberTest) getReq() interface{} {
	return t.Req
}

func (t DeleteMemberTest) getExpCode() interface{} {
	return t.ExpCode
}

func (t DeleteMemberTest) getExpResp() interface{} {
	return t.ExpResp
}
