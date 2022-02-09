package test

import (
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

// AssertCase Run the testCase and assert whether the result is equal to the expected value.
func AssertCase(t *testing.T, router *gin.Engine, method string, pathPrefix string, relativePath string, testCase BaseTest) {
	w := CallApi(router, method, pathPrefix, relativePath, testCase.getReq())
	assert.Equal(t, testCase.getExpCode(), w.Code)
	expResp, _ := json.Marshal(testCase.getExpResp())
	assert.Equal(t, string(expResp), w.Body.String())
}

func AssertBenchmarkCase(b *testing.B, router *gin.Engine, method string, pathPrefix string, relativePath string, testCase BaseTest) {
	w := CallApi(router, method, pathPrefix, relativePath, testCase.getReq())
	assert.Equal(b, testCase.getExpCode(), w.Code)
	expResp, _ := json.Marshal(testCase.getExpResp())
	assert.Equal(b, string(expResp), w.Body.String())
}
func AssertCaseCookie(t *testing.T, router *gin.Engine, w *httptest.ResponseRecorder, method string, pathPrefix string, relativePath string, testCase BaseTest) {
	//w := CallApi(router, method, pathPrefix, relativePath, testCase.getReq())
	//assert.Equal(t, testCase.getExpCode(), w.Code)
	expResp, _ := json.Marshal(testCase.getExpResp())
	assert.Equal(t, string(expResp), w.Body.String())
}

// UserTest with author check

func CallApiAndAddCookie(router *gin.Engine, method string, pathPrefix string, relativePath string, reqData interface{}) (w *httptest.ResponseRecorder) {
	cookie := http.Cookie{
		Name:       "camp-session",
		Value:      "1",
		Path:       "/",
		Domain:     "180.184.74.137",
		RawExpires: "",
		MaxAge:     0,
		Secure:     false,
		HttpOnly:   true,
	}
	w = httptest.NewRecorder()
	body, _ := json.Marshal(reqData)
	req, _ := http.NewRequest(method, pathPrefix+relativePath, strings.NewReader(string(body)))
	req.AddCookie(&cookie)
	router.ServeHTTP(w, req)
	return w
}

func AssertCaseForCreate(t *testing.T, router *gin.Engine, method string, pathPrefix string, relativePath string, testCase BaseTest) {
	w := CallApiAndAddCookie(router, method, pathPrefix, relativePath, testCase.getReq())
	assert.Equal(t, testCase.getExpCode(), w.Code)
	expResp, _ := json.Marshal(testCase.getExpResp())
	assert.Equal(t, string(expResp), w.Body.String())
}

func AssertBenchmarkCaseForCreate(b *testing.B, router *gin.Engine, method string, pathPrefix string, relativePath string, testCase BaseTest) {
	w := CallApiAndAddCookie(router, method, pathPrefix, relativePath, testCase.getReq())
	assert.Equal(b, testCase.getExpCode(), w.Code)
	expResp, _ := json.Marshal(testCase.getExpResp())
	assert.Equal(b, string(expResp), w.Body.String())
}
