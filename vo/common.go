package vo

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"strings"
)

func CallApi(router *gin.Engine, method string, pathPrefix string, relativePath string, reqData interface{}) (w *httptest.ResponseRecorder) {
	w = httptest.NewRecorder()
	body, _ := json.Marshal(reqData)
	req, _ := http.NewRequest(method, pathPrefix+relativePath, strings.NewReader(string(body)))
	router.ServeHTTP(w, req)
	return w
}
