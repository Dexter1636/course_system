package main

import (
	"course_system/common"
	"course_system/test"
	"course_system/test/cases"
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
	initDataForCourseBooking()

	for _, tc := range cases.BookCourseCases {
		test.AssertCase(t, router, "POST", pathPrefix, "/student/book_course", tc)
	}
}

func BenchmarkBookCourseRoute(b *testing.B) {
	b.Cleanup(cleanup)
	initDataForCourseBooking()

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
