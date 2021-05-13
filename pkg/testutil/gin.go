package testutil

import (
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
)

func TestGinHTTPResponse(
	r *gin.Engine,
	req *http.Request,
	testFn func(w *httptest.ResponseRecorder),
) {
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	testFn(w)
}
