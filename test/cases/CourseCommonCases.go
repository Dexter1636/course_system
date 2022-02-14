package cases

import (
	"course_system/test"
	"course_system/vo"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
)

var CreateCourseCases = []test.CreateCourseTest{
	{
		Req: vo.CreateCourseRequest{
			Name: "Introduction to C++",
			Cap:  120,
		},
		ExpCode: http.StatusOK,
		ExpResp: vo.CreateCourseResponse{
			Code: vo.OK,
			Data: struct {
				CourseID string
			}{CourseID: "1"},
		},
	},
	{
		Req: vo.CreateCourseRequest{
			Name: "Introduction to Java",
			Cap:  140,
		},
		ExpCode: http.StatusOK,
		ExpResp: vo.CreateCourseResponse{
			Code: vo.OK,
			Data: struct {
				CourseID string
			}{CourseID: "2"},
		},
	},
}

func GenerateCreateCourseCase(i int) (tc test.CreateCourseTest) {
	tc = test.CreateCourseTest{
		Req: vo.CreateCourseRequest{
			Name: fmt.Sprintf("Test Course %d", i),
			Cap:  rand.Intn(1000),
		},
		ExpCode: http.StatusOK,
		ExpResp: vo.CreateCourseResponse{
			Code: vo.OK,
			Data: struct {
				CourseID string
			}{CourseID: strconv.Itoa(i + 1)},
		},
	}
	return tc
}

var GetCourseCases = []test.GetCourseTest{
	{
		Req:     vo.GetCourseRequest{CourseID: "1"},
		ExpCode: http.StatusOK,
		ExpResp: vo.GetCourseResponse{
			Code: vo.CourseNotExisted,
			Data: vo.TCourse{},
		},
	},
}

func GenerateGetCourseCase(i int) (tc test.GetCourseTest) {
	tc = test.GetCourseTest{
		Req:     vo.GetCourseRequest{CourseID: strconv.FormatInt(rand.Int63n(1000), 10)},
		ExpCode: http.StatusOK,
		ExpResp: vo.GetCourseResponse{
			Code: vo.CourseNotExisted,
			Data: vo.TCourse{},
		},
	}
	return tc
}
