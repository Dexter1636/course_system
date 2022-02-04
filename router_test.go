package main

import (
	"course_system/common"
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
	common.GetDB().Raw("TRUNCATE TABLE user")
	common.GetDB().Raw("TRUNCATE TABLE course")
	common.GetDB().Raw("TRUNCATE TABLE sc")
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

	assert.Equal(t, 200, w.Code)
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
	w := httptest.NewRecorder()
	body, _ := json.Marshal(vo.CreateCourseRequest{
		Name: "Introduction to C++",
		Cap:  120,
	})
	req, _ := http.NewRequest("POST", pathPrefix+"/course/create", strings.NewReader(string(body)))
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	var resp vo.GetCourseResponse
	if err := json.Unmarshal([]byte(w.Body.String()), &resp); err != nil {
		panic(err.Error())
	}
	assert.Equal(t, vo.CreateCourseResponse{
		Code: vo.OK,
		Data: struct {
			CourseID string
		}{CourseID: "1"},
	}, resp)

	t.Cleanup(cleanup)
}

func TestGetCourseRoute(t *testing.T) {
	w := httptest.NewRecorder()
	body, _ := json.Marshal(vo.GetCourseRequest{CourseID: "1"})
	req, _ := http.NewRequest("GET", pathPrefix+"/course/get", strings.NewReader(string(body)))
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	var resp vo.GetCourseResponse
	if err := json.Unmarshal([]byte(w.Body.String()), &resp); err != nil {
		panic(err.Error())
	}
	assert.Equal(t, vo.GetCourseResponse{
		Code: vo.CourseNotExisted,
		Data: vo.TCourse{},
	}, resp)

	t.Cleanup(cleanup)
}

func BenchmarkGetCourseRoute(b *testing.B) {
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		body, _ := json.Marshal(vo.GetCourseRequest{CourseID: strconv.FormatInt(rand.Int63n(15), 10)})
		req, _ := http.NewRequest("GET", pathPrefix+"/course/get", strings.NewReader(string(body)))
		router.ServeHTTP(w, req)
	}
}
