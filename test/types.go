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
