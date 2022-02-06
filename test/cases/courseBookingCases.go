package cases

import (
	"course_system/test"
	"course_system/vo"
	"net/http"
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
