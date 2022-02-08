package main

import (
	"context"
	"course_system/common"
	"course_system/test"
	"course_system/test/cases"
	"course_system/test/data"
	"fmt"
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
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

//=============== User ========
func TestGetMemberListRoute(t *testing.T) {

	t.Cleanup(cleanup)

	data.InitDataForUser()
	data.InitDataForUserOther()

	for _, tc := range cases.GetMemberListCases {
		test.AssertCase(t, router, "GET", pathPrefix, "/member/list", tc)
	}
}

func BenchmarkGetMemberListRoute(b *testing.B) {
	b.Cleanup(cleanup)

	data.InitDataForUser()
	data.InitDataForUserOther()

	for i := 0; i < b.N; i++ {
		test.AssertBenchmarkCase(b, router, "GET", pathPrefix, "/member/list", cases.GenerateGetMemberListCase(i))
	}
}

func TestUpdateMemberRoute(t *testing.T) {
	t.Cleanup(cleanup)

	data.InitDataForUser()
	data.InitDataForUserOther()

	for _, tc := range cases.UpdateMemberCases {
		test.AssertCase(t, router, "POST", pathPrefix, "/member/update", tc)
	}
}

func BenchmarkUpdateMemberRoute(b *testing.B) {
	b.Cleanup(cleanup)

	data.InitDataForUser()
	data.InitDataForUserOther()

	for i := 0; i < b.N; i++ {
		test.AssertBenchmarkCase(b, router, "POST", pathPrefix, "/member/update", cases.GenerateUpdateMemberCase(i))
	}
}

func TestDeleteMemberRoute(t *testing.T) {
	t.Cleanup(cleanup)

	data.InitDataForUser()
	data.InitDataForUserOther()

	for _, tc := range cases.DeleteMemberCases {
		test.AssertCase(t, router, "POST", pathPrefix, "/member/delete", tc)
	}
}

func BenchmarkDeleteMemberRoute(b *testing.B) {
	b.Cleanup(cleanup)

	data.InitDataForUser()
	data.InitDataForUserOther()

	for i := 0; i < b.N; i++ {
		test.AssertBenchmarkCase(b, router, "POST", pathPrefix, "/member/delete", cases.GenerateDeleteMemberCase(i))
	}
}

//================Course Schedule=========================
func TestBindCourseRoute(t *testing.T) {
	t.Cleanup(cleanup)

	for _, tc := range cases.BindCoruseCases {
		test.AssertCase(t, router, "POST", pathPrefix, "/teacher/bind_course", tc)
	}
}
func BenchmarkBindRoute(b *testing.B) {
	b.Cleanup(cleanup)
	data.InitDataForCourseCommon()
	for i := 0; i < b.N; i++ {
		test.AssertBenchmarkCase(b, router, "POST", pathPrefix, "/teacher/bind_course", cases.GenerateBingCase(i))
	}
}
func TestUnBindCourseRoute(t *testing.T) {
	t.Cleanup(cleanup)

	for _, tc := range cases.UnbindCourseCases {
		test.AssertCase(t, router, "POST", pathPrefix, "/teacher/unbind_course", tc)
	}
}
func BenchmarkUnbindRoute(b *testing.B) {
	b.Cleanup(cleanup)
	data.InitDataForUnbing()
	for i := 0; i < b.N; i++ {
		test.AssertBenchmarkCase(b, router, "POST", pathPrefix, "/teacher/unbind_course", cases.GenerateUnbingCase(i))
	}
}

func TGetCourseRoute(t *testing.T) {
	t.Cleanup(cleanup)

	for _, tc := range cases.GetTCourseCases {
		test.AssertCase(t, router, "GET", pathPrefix, "/teacher/get_course", tc)
	}
}

func BenchmarkTGetCourseRoute(b *testing.B) {
	b.Cleanup(cleanup)
	data.InitDataForUnbing()
	for i := 0; i < b.N; i++ {
		test.AssertBenchmarkCase(b, router, "GET", pathPrefix, "/teacher/get_course", cases.GenerateTGetcourse(i))
	}
}

func TestScheduleCourseRoute(t *testing.T) {
	t.Cleanup(cleanup)

	for _, tc := range cases.ScheduleCases {
		test.AssertCase(t, router, "POST", pathPrefix, "/course/schedule", tc)
	}
}
func BenchmarkScheduleRouter(b *testing.B) {
	b.Cleanup(cleanup)

	for i := 0; i < b.N; i++ {
		test.AssertBenchmarkCase(b, router, "POST", pathPrefix, "/course/schedule", cases.GenerateSchedule(i))
	}
}
