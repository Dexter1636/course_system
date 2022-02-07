package test

import (
	"course_system/vo"
)

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
type ScheduleTest struct {
	Req     vo.ScheduleCourseRequest
	ExpCode int
	ExpResp vo.ScheduleCourseResponse
}
