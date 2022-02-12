package cases

import (
	"course_system/test"
	"course_system/vo"
	"math/rand"
	"net/http"
	"strconv"
)

var BindCoruseCases = []test.BindCourseTest{
	{
		Req: vo.BindCourseRequest{
			CourseID:  strconv.FormatInt(rand.Int63n(1000), 10),
			TeacherID: strconv.FormatInt(rand.Int63n(1000), 10),
		},
		ExpCode: http.StatusOK,
		ExpResp: vo.BindCourseResponse{
			Code: vo.CourseNotExisted,
		},
	},
}

func GenerateBingCase(i int) (tc test.BindCourseTest) {
	tc = test.BindCourseTest{
		Req: vo.BindCourseRequest{
			CourseID:  strconv.FormatInt(int64(i%5+1), 10),
			TeacherID: strconv.FormatInt(rand.Int63n(1000), 10),
		},
		ExpCode: http.StatusOK,
		ExpResp: vo.BindCourseResponse{
			vo.OK,
		},
	}
	if i%5 == 4 {
		tc.ExpResp.Code = vo.CourseNotExisted
	} else if i >= 5 {
		tc.ExpResp.Code = vo.CourseHasBound
	}
	return tc
}

var UnbindCourseCases = []test.UnBindCourseTest{
	{
		Req: vo.UnbindCourseRequest{
			CourseID:  strconv.FormatInt(rand.Int63n(1000), 10),
			TeacherID: strconv.FormatInt(rand.Int63n(1000), 10),
		},
		ExpCode: http.StatusOK,
		ExpResp: vo.UnbindCourseResponse{
			Code: vo.CourseNotExisted,
		},
	},
}

func GenerateUnbingCase(i int) (tc test.UnBindCourseTest) {
	tc = test.UnBindCourseTest{
		Req: vo.UnbindCourseRequest{
			CourseID:  strconv.FormatInt(int64(i%5+1), 10),
			TeacherID: strconv.FormatInt(rand.Int63n(1000), 10),
		},
		ExpCode: http.StatusOK,
		ExpResp: vo.UnbindCourseResponse{
			vo.OK,
		},
	}
	if i%5 == 4 {
		tc.ExpResp.Code = vo.CourseNotExisted
	} else if i >= 5 {
		tc.ExpResp.Code = vo.CourseNotBind
	}
	return tc
}

var GetTCourseCases = []test.TGetCourseTests{
	{
		Req: vo.GetTeacherCourseRequest{
			TeacherID: strconv.FormatInt(rand.Int63n(1000), 10),
		},
		ExpCode: http.StatusOK,
		ExpResp: vo.GetTeacherCourseResponse{
			Code: vo.OK,
			Data: struct{ CourseList []*vo.TCourse }{CourseList: nil},
		},
	},
}

func GenerateTGetcourse(i int) (tc test.TGetCourseTests) {
	var a []*vo.TCourse = make([]*vo.TCourse, 2)
	a[0] = new(vo.TCourse)
	a[0].Name = "test2"
	a[0].TeacherID = "893"
	a[0].CourseID = "2"
	a[1] = new(vo.TCourse)
	a[1].Name = "test4"
	a[1].TeacherID = "893"
	a[1].CourseID = "4"
	var b []*vo.TCourse = make([]*vo.TCourse, 2)
	b[0] = new(vo.TCourse)
	b[0].Name = "test1"
	b[0].TeacherID = "810"
	b[0].CourseID = "1"
	b[1] = new(vo.TCourse)
	b[1].Name = "test3"
	b[1].TeacherID = "810"
	b[1].CourseID = "3"
	tc = test.TGetCourseTests{
		Req: vo.GetTeacherCourseRequest{
			TeacherID: strconv.FormatInt(int64(893), 10),
		},
		ExpCode: http.StatusOK,
		ExpResp: vo.GetTeacherCourseResponse{
			Code: vo.OK,
			Data: struct{ CourseList []*vo.TCourse }{CourseList: a},
		},
	}
	if i%2 == 1 {
		tc.Req.TeacherID = "810"
		tc.ExpResp.Data = struct{ CourseList []*vo.TCourse }{CourseList: b}
	}
	return tc
}

var ScheduleCases = []test.ScheduleTest{
	{
		Req: vo.ScheduleCourseRequest{
			TeacherCourseRelationShip: map[string][]string{
				"TNOK": {"893", "810"},
				"DB":   {"893", "114"},
				"TDN":  {"810", "514"},
				"MUR":  {"114"},
			},
		},
		ExpCode: http.StatusOK,
		ExpResp: vo.ScheduleCourseResponse{
			Code: vo.OK,
			Data: map[string]string{
				"MUR":  "114",
				"DB":   "893",
				"TNOK": "810",
				"TDN":  "514",
			},
		},
	},
}

func GenerateSchedule(i int) (tc test.ScheduleTest) {
	tc = test.ScheduleTest{
		Req: vo.ScheduleCourseRequest{
			TeacherCourseRelationShip: map[string][]string{
				"TNOK": {"893", "810"},
				"DB":   {"893", "114"},
				"TDN":  {"810", "514"},
				"MUR":  {"114"},
			},
		},
		ExpCode: http.StatusOK,
		ExpResp: vo.ScheduleCourseResponse{
			Code: vo.OK,
			Data: map[string]string{
				"MUR":  "114",
				"DB":   "893",
				"TNOK": "810",
				"TDN":  "514",
			},
		},
	}
	return tc
}
