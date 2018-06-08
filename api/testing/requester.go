package testing

import (
	"bytes"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	. "github.com/onsi/gomega"
)

type Requester struct {
	URL     string
	Handler gin.HandlerFunc
}

func (r *Requester) Get() *httptest.ResponseRecorder {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.GET(r.URL, r.Handler)

	recorder := httptest.NewRecorder()

	request, err := http.NewRequest("GET", r.URL, nil)

	Expect(err).To(BeNil())

	router.ServeHTTP(recorder, request)

	return recorder
}

func (r *Requester) Post(data []byte) *httptest.ResponseRecorder {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.POST(r.URL, r.Handler)

	recorder := httptest.NewRecorder()

	request, err := http.NewRequest("POST", r.URL, bytes.NewReader(data))

	Expect(err).To(BeNil())

	router.ServeHTTP(recorder, request)

	return recorder
}
