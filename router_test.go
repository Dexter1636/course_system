package main

import (
	"course_system/common"
	"course_system/test"
	"course_system/test/cases"
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

// ==== DB Data Init Functions ====

func initDataForCourseCommon() {
	// insert students
	common.GetDB().Exec("INSERT INTO user(user_name, nick_name, password, role_id, enabled) VALUES ('Amy Wong', 'Amy', '123456', '2', 1)")
	common.GetDB().Exec("INSERT INTO user(user_name, nick_name, password, role_id, enabled) VALUES ('Dexter Peng', 'Dexter', '123456', '2', 1)")
	common.GetDB().Exec("INSERT INTO user(user_name, nick_name, password, role_id, enabled) VALUES ('San Zhang', 'San', '123456', '2', 1)")
	common.GetDB().Exec("INSERT INTO user(user_name, nick_name, password, role_id, enabled) VALUES ('Si Li', 'Si', '123456', '2', 1)")

	// insert courses
	common.GetDB().Exec("INSERT INTO course(name, avail, cap) VALUES ('test1', 1, 1)")
	common.GetDB().Exec("INSERT INTO course(name, avail, cap) VALUES ('test2', 3, 3)")
	common.GetDB().Exec("INSERT INTO course(name, avail, cap) VALUES ('test3', 0, 100)")
	common.GetDB().Exec("INSERT INTO course(name, avail, cap) VALUES ('test4', 100, 100)")
}

func initDataForCourseBooking() {
	initDataForCourseCommon()
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

	tests := []test.CreateCourseTest{
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

	for _, tc := range tests {
		w := httptest.NewRecorder()
		body, _ := json.Marshal(tc.Req)
		req, _ := http.NewRequest("POST", pathPrefix+"/course/create", strings.NewReader(string(body)))
		router.ServeHTTP(w, req)

		assert.Equal(t, tc.ExpCode, w.Code)
		var resp vo.CreateCourseResponse
		if err := json.Unmarshal([]byte(w.Body.String()), &resp); err != nil {
			panic(err.Error())
		}
		assert.Equal(t, tc.ExpResp, resp)
	}
}

func BenchmarkCreateCourseRoute(b *testing.B) {
	b.Cleanup(cleanup)

	for i := 0; i < b.N; i++ {
		tc := test.CreateCourseTest{
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

		w := httptest.NewRecorder()
		body, _ := json.Marshal(tc.Req)
		req, _ := http.NewRequest("POST", pathPrefix+"/course/create", strings.NewReader(string(body)))
		router.ServeHTTP(w, req)

		assert.Equal(b, tc.ExpCode, w.Code)
		var resp vo.CreateCourseResponse
		if err := json.Unmarshal([]byte(w.Body.String()), &resp); err != nil {
			panic(err.Error())
		}
		assert.Equal(b, tc.ExpResp, resp)
	}
}

func TestGetCourseRoute(t *testing.T) {
	t.Cleanup(cleanup)

	tests := []test.GetCourseTest{
		{
			Req:     vo.GetCourseRequest{CourseID: "1"},
			ExpCode: http.StatusOK,
			ExpResp: vo.GetCourseResponse{
				Code: vo.CourseNotExisted,
				Data: vo.TCourse{},
			},
		},
	}

	for _, tc := range tests {
		w := httptest.NewRecorder()
		body, _ := json.Marshal(tc.Req)
		req, _ := http.NewRequest("GET", pathPrefix+"/course/get", strings.NewReader(string(body)))
		router.ServeHTTP(w, req)

		assert.Equal(t, tc.ExpCode, w.Code)
		var resp vo.GetCourseResponse
		if err := json.Unmarshal([]byte(w.Body.String()), &resp); err != nil {
			panic(err.Error())
		}
		assert.Equal(t, tc.ExpResp, resp)
	}
}

func BenchmarkGetCourseRoute(b *testing.B) {
	b.Cleanup(cleanup)

	for i := 0; i < b.N; i++ {
		tc := test.GetCourseTest{
			Req:     vo.GetCourseRequest{CourseID: strconv.FormatInt(rand.Int63n(1000), 10)},
			ExpCode: http.StatusOK,
			ExpResp: vo.GetCourseResponse{
				Code: vo.CourseNotExisted,
				Data: vo.TCourse{},
			},
		}

		w := httptest.NewRecorder()
		body, _ := json.Marshal(tc.Req)
		req, _ := http.NewRequest("GET", pathPrefix+"/course/get", strings.NewReader(string(body)))
		router.ServeHTTP(w, req)

		assert.Equal(b, tc.ExpCode, w.Code)
		var resp vo.GetCourseResponse
		if err := json.Unmarshal([]byte(w.Body.String()), &resp); err != nil {
			panic(err.Error())
		}
		assert.Equal(b, tc.ExpResp, resp)
	}
}

// ======== CourseBooking ========

func TestBookCourseRoute(t *testing.T) {
	t.Cleanup(cleanup)
	initDataForCourseBooking()

	for _, tc := range cases.BookCourseCases {
		test.AssertCase(t, router, "POST", pathPrefix, "/student/book_course", tc)
	}
}

func BenchmarkBookCourseRoute(b *testing.B) {
	b.Cleanup(cleanup)
	initDataForCourseBooking()

	for i := 0; i < b.N; i++ {
		test.CallApi(router, "POST", pathPrefix, "/student/book_course", cases.GenerateBookCourseReq())
	}
}

func TestGetStudentCourseRoute(t *testing.T) {
	t.Cleanup(cleanup)

	tests := []test.GetStudentCourseTest{
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

	for _, tc := range tests {
		test.AssertCase(t, router, "GET", pathPrefix, "/student/course", tc)
	}
}

func BenchmarkGetStudentCourseRoute(b *testing.B) {
	b.Cleanup(cleanup)
}
