package cases

import (
	"course_system/test"
	"course_system/vo"
	"math/rand"
	"net/http"
	"strconv"
)

var BookCourseCases = []test.BookCourseTest{
	{
		Req: vo.BookCourseRequest{
			StudentID: "1",
			CourseID:  "1",
		},
		ExpCode: http.StatusOK,
		ExpResp: vo.BookCourseResponse{Code: vo.OK},
	},
	{
		Req: vo.BookCourseRequest{
			StudentID: "2",
			CourseID:  "1",
		},
		ExpCode: http.StatusOK,
		ExpResp: vo.BookCourseResponse{Code: vo.CourseNotAvailable},
	},
	{
		Req: vo.BookCourseRequest{
			StudentID: "2",
			CourseID:  "3",
		},
		ExpCode: http.StatusOK,
		ExpResp: vo.BookCourseResponse{Code: vo.CourseNotAvailable},
	},
	{
		Req: vo.BookCourseRequest{
			StudentID: "2",
			CourseID:  "2",
		},
		ExpCode: http.StatusOK,
		ExpResp: vo.BookCourseResponse{Code: vo.OK},
	},
}

func GenerateBookCourseReq() (tc vo.BookCourseRequest) {
	tc = vo.BookCourseRequest{
		StudentID: strconv.FormatInt(rand.Int63n(10), 10),
		CourseID:  strconv.FormatInt(rand.Int63n(10), 10),
	}
	return tc
}

var GetStudentCourseCases = []test.GetStudentCourseTest{
	{
		Req:     vo.GetStudentCourseRequest{StudentID: "1"},
		ExpCode: http.StatusOK,
		ExpResp: vo.GetStudentCourseResponse{
			Code: vo.StudentNotExisted,
			Data: struct {
				CourseList []vo.TCourse
			}{CourseList: []vo.TCourse{}},
		},
	},
}
