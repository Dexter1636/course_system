package test

import (
	"course_system/vo"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func CallApi(router *gin.Engine, method string, pathPrefix string, relativePath string, reqData interface{}) (w *httptest.ResponseRecorder) {
	w = httptest.NewRecorder()
	body, _ := json.Marshal(reqData)
	req, _ := http.NewRequest(method, pathPrefix+relativePath, strings.NewReader(string(body)))
	router.ServeHTTP(w, req)
	return w
}

func AssertCase(t *testing.T, router *gin.Engine, method string, pathPrefix string, relativePath string, testCase BaseTest) {
	w := CallApi(router, method, pathPrefix, relativePath, testCase.getReq())
	assert.Equal(t, testCase.getExpCode(), w.Code)
	var resp vo.GetStudentCourseResponse
	if err := json.Unmarshal([]byte(w.Body.String()), &resp); err != nil {
		panic(err.Error())
	}
	assert.Equal(t, testCase.getExpResp(), resp)
}
