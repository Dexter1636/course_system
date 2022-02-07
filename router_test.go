package main

import (
	"context"
	"course_system/common"
	"course_system/test"
	"course_system/test/cases"
	"course_system/test/data"
	"course_system/vo"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var router *gin.Engine
var pathPrefix string

// before all
func setup() {
	common.InitConfig("test")
	common.InitDb()
	common.InitRdb(context.Background())
	router = RegisterRouter()
	pathPrefix = "/api/v1"
	rand.Seed(10)
	cleanup()
}

// after all
func teardown() {
	cleanup()
}

// after each test
func cleanup() {
	fmt.Println("==cleanup==")
	common.GetDB().Exec("TRUNCATE TABLE user")
	common.GetDB().Exec("TRUNCATE TABLE course")
	common.GetDB().Exec("TRUNCATE TABLE sc")
}

func TestMain(m *testing.M) {
	setup()
	fmt.Println("Test begins....")
	code := m.Run()
	teardown()
	os.Exit(code)
}

// ======== Ping ========

func TestPingRoute(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", pathPrefix+"/ping", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "{\"message\":\"pong\"}", w.Body.String())
}

func BenchmarkPingRoute(b *testing.B) {
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", pathPrefix+"/ping", nil)
		router.ServeHTTP(w, req)
	}
}

// ======== CourseCommon ========

func TestCreateCourseRoute(t *testing.T) {
	t.Cleanup(cleanup)

	for _, tc := range cases.CreateCourseCases {
		test.AssertCase(t, router, "POST", pathPrefix, "/course/create", tc)
	}
}

func BenchmarkCreateCourseRoute(b *testing.B) {
	b.Cleanup(cleanup)

	for i := 0; i < b.N; i++ {
		test.AssertBenchmarkCase(b, router, "POST", pathPrefix, "/course/create", cases.GenerateCreateCourseCase(i))
	}
}

func TestGetCourseRoute(t *testing.T) {
	t.Cleanup(cleanup)

	for _, tc := range cases.GetCourseCases {
		test.AssertCase(t, router, "GET", pathPrefix, "/course/get", tc)
	}
}

func BenchmarkGetCourseRoute(b *testing.B) {
	b.Cleanup(cleanup)

	for i := 0; i < b.N; i++ {
		test.AssertBenchmarkCase(b, router, "GET", pathPrefix, "/course/get", cases.GenerateGetCourseCase(i))
	}
}

// ======== CourseBooking ========

func TestBookCourseRoute(t *testing.T) {
	t.Cleanup(cleanup)
	data.InitDataForCourseBooking()

	for _, tc := range cases.BookCourseCases {
		test.AssertCase(t, router, "POST", pathPrefix, "/student/book_course", tc)
	}
}

func BenchmarkBookCourseRoute(b *testing.B) {
	b.Cleanup(cleanup)
	data.InitDataForCourseBooking()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			test.CallApi(router, "POST", pathPrefix, "/student/book_course", cases.GenerateBookCourseReq())
		}
	})
}

func TestGetStudentCourseRoute(t *testing.T) {
	t.Cleanup(cleanup)

	for _, tc := range cases.GetStudentCourseCases {
		test.AssertCase(t, router, "GET", pathPrefix, "/student/course", tc)
	}
}

func BenchmarkGetStudentCourseRoute(b *testing.B) {
	b.Cleanup(cleanup)
}

// ======== User ========(Create && Get)
//测试前需将UserController中的“权限检查”部分注释掉

func TestCreateMemberRoute(t *testing.T) {
	t.Cleanup(cleanup)

	for _, tc := range cases.CreateMemberCases {
		test.AssertCase(t, router, "POST", pathPrefix, "/member/create", tc)
	}
}

func BenchmarkCreateMemberRoute(b *testing.B) {
	b.Cleanup(cleanup)

	for i := 0; i < b.N; i++ {
		test.AssertBenchmarkCase(b, router, "POST", pathPrefix, "/member/create", cases.GenerateCreateMemberCase(i))
	}
}

func TestGetMemberRoute(t *testing.T) {
	t.Cleanup(cleanup)

	data.InitDataForUser()

	for _, tc := range cases.GetMemberCases {
		test.AssertCase(t, router, "GET", pathPrefix, "/member", tc)
	}
}

func BenchmarkGetMemberRoute(b *testing.B) {
	b.Cleanup(cleanup)

	data.InitDataForUser()

	for i := 0; i < b.N; i++ {
		test.AssertBenchmarkCase(b, router, "GET", pathPrefix, "/member", cases.GenerateGetMemberCase(i))
	}
}
func TestBindRoute(t *testing.T) {
	tests := []test.BindCourseTest{
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
	for _, tc := range tests {
		w := httptest.NewRecorder()
		body, _ := json.Marshal(tc.Req)
		req, _ := http.NewRequest("POST", pathPrefix+"/teacher/bind_course", strings.NewReader(string(body)))
		router.ServeHTTP(w, req)
		assert.Equal(t, tc.ExpCode, w.Code)
		var resp vo.BindCourseResponse
		if err := json.Unmarshal([]byte(w.Body.String()), &resp); err != nil {
			panic(err.Error())
		}
		assert.Equal(t, tc.ExpResp, resp)
	}

}
func BenchmarkTestBindRoute(b *testing.B) {
	for i := 0; i < b.N; i++ {
		tests := []test.BindCourseTest{
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
		for _, tc := range tests {
			w := httptest.NewRecorder()
			body, _ := json.Marshal(tc.Req)
			req, _ := http.NewRequest("POST", pathPrefix+"/teacher/bind_course", strings.NewReader(string(body)))
			router.ServeHTTP(w, req)
			assert.Equal(b, tc.ExpCode, w.Code)
			var resp vo.BindCourseResponse
			if err := json.Unmarshal([]byte(w.Body.String()), &resp); err != nil {
				panic(err.Error())
			}
			assert.Equal(b, tc.ExpResp, resp)
		}
	}
}
func TestUnBindRoute(t *testing.T) {
	tests := []test.UnBindCourseTest{
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
	for _, tc := range tests {
		w := httptest.NewRecorder()
		body, _ := json.Marshal(tc.Req)
		req, _ := http.NewRequest("POST", pathPrefix+"/teacher/unbind_course", strings.NewReader(string(body)))
		router.ServeHTTP(w, req)
		assert.Equal(t, tc.ExpCode, w.Code)
		var resp vo.UnbindCourseResponse
		if err := json.Unmarshal([]byte(w.Body.String()), &resp); err != nil {
			panic(err.Error())
		}
		assert.Equal(t, tc.ExpResp, resp)
	}

}
func BenchmarkTestUnBindRoute(b *testing.B) {
	for i := 0; i < b.N; i++ {
		tests := []test.UnBindCourseTest{
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
		for _, tc := range tests {
			w := httptest.NewRecorder()
			body, _ := json.Marshal(tc.Req)
			req, _ := http.NewRequest("POST", pathPrefix+"/teacher/unbind_course", strings.NewReader(string(body)))
			router.ServeHTTP(w, req)
			assert.Equal(b, tc.ExpCode, w.Code)
			var resp vo.UnbindCourseResponse
			if err := json.Unmarshal([]byte(w.Body.String()), &resp); err != nil {
				panic(err.Error())
			}
			assert.Equal(b, tc.ExpResp, resp)
		}
	}
}
func TestGetRoute(t *testing.T) {
	tests := []test.GetCourseTest{
		{
			Req: vo.GetCourseRequest{
				CourseID: strconv.FormatInt(rand.Int63n(1000), 10),
			},
			ExpCode: http.StatusOK,
			ExpResp: vo.GetCourseResponse{
				Code: vo.OK,
			},
		},
	}
	for _, tc := range tests {
		w := httptest.NewRecorder()
		body, _ := json.Marshal(tc.Req)
		req, _ := http.NewRequest("GET", pathPrefix+"/teacher/get_course", strings.NewReader(string(body)))
		router.ServeHTTP(w, req)
		assert.Equal(t, tc.ExpCode, w.Code)
		var resp vo.GetCourseResponse
		if err := json.Unmarshal([]byte(w.Body.String()), &resp); err != nil {
			panic(err.Error())
		}
		assert.Equal(t, tc.ExpResp, resp)
	}

}
func BenchmarkTestGetRoute(b *testing.B) {
	for i := 0; i < b.N; i++ {
		tests := []test.GetCourseTest{
			{
				Req: vo.GetCourseRequest{
					CourseID: strconv.FormatInt(rand.Int63n(1000), 10),
				},
				ExpCode: http.StatusOK,
				ExpResp: vo.GetCourseResponse{
					Code: vo.OK,
				},
			},
		}
		for _, tc := range tests {
			w := httptest.NewRecorder()
			body, _ := json.Marshal(tc.Req)
			req, _ := http.NewRequest("GET", pathPrefix+"/teacher/get_course", strings.NewReader(string(body)))
			router.ServeHTTP(w, req)
			assert.Equal(b, tc.ExpCode, w.Code)
			var resp vo.GetCourseResponse
			if err := json.Unmarshal([]byte(w.Body.String()), &resp); err != nil {
				panic(err.Error())
			}
			assert.Equal(b, tc.ExpResp, resp)
		}
	}
}
func TestScheduleRoute(t *testing.T) {

	var data map[string][]string = make(map[string][]string)
	var ans map[string]string = make(map[string]string)
	var a = []string{"893", "810"}
	var b = []string{"893", "114"}
	var c = []string{"810", "514"}
	var d = []string{"114"}

	data["TNOK"] = a
	data["DB"] = b
	data["TDN"] = c
	data["MUR"] = d

	ans["MUR"] = "114"
	ans["DB"] = "893"
	ans["TNOK"] = "810"
	ans["TDN"] = "514"

	tests := []test.ScheduleTest{
		{
			Req: vo.ScheduleCourseRequest{
				TeacherCourseRelationShip: data,
			},
			ExpCode: http.StatusOK,
			ExpResp: vo.ScheduleCourseResponse{
				Code: vo.OK,
				Data: ans,
			},
		},
	}
	for _, tc := range tests {
		w := httptest.NewRecorder()
		body, _ := json.Marshal(tc.Req)
		req, _ := http.NewRequest("POST", pathPrefix+"/course/schedule", strings.NewReader(string(body)))
		router.ServeHTTP(w, req)
		assert.Equal(t, tc.ExpCode, w.Code)
		var resp vo.ScheduleCourseResponse
		if err := json.Unmarshal([]byte(w.Body.String()), &resp); err != nil {
			panic(err.Error())
		}
		assert.Equal(t, tc.ExpResp, resp)
	}
}
func BenchmarkTestScheduleRoute(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var data map[string][]string = make(map[string][]string)
		var ans map[string]string = make(map[string]string)
		var a = []string{"893", "810"}
		var B = []string{"893", "114"}
		var c = []string{"810", "514"}
		var d = []string{"114"}

		data["TNOK"] = a
		data["DB"] = B
		data["TDN"] = c
		data["MUR"] = d

		ans["MUR"] = "114"
		ans["DB"] = "893"
		ans["TNOK"] = "810"
		ans["TDN"] = "514"

		tests := []test.ScheduleTest{
			{
				Req: vo.ScheduleCourseRequest{
					TeacherCourseRelationShip: data,
				},
				ExpCode: http.StatusOK,
				ExpResp: vo.ScheduleCourseResponse{
					Code: vo.OK,
					Data: ans,
				},
			},
		}
		for _, tc := range tests {
			w := httptest.NewRecorder()
			body, _ := json.Marshal(tc.Req)
			req, _ := http.NewRequest("POST", pathPrefix+"/course/schedule", strings.NewReader(string(body)))
			router.ServeHTTP(w, req)
			assert.Equal(b, tc.ExpCode, w.Code)
			var resp vo.ScheduleCourseResponse
			if err := json.Unmarshal([]byte(w.Body.String()), &resp); err != nil {
				panic(err.Error())
			}
			assert.Equal(b, tc.ExpResp, resp)
		}
	}
}
