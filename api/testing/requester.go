package testing

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/gin-gonic/gin"
	. "github.com/onsi/gomega"
)

type Requester struct {
	URL     string
	Handler gin.HandlerFunc
}

type Params map[string]string

func (r *Requester) Get(params Params, query Params) *httptest.ResponseRecorder {
	router := r.server()

	url := r.parseParams(r.URL, params)
	url = r.parseQuery(url, query)

	recorder := httptest.NewRecorder()
	request, err := http.NewRequest("GET", url, nil)

	Expect(err).To(BeNil())

	router.ServeHTTP(recorder, request)

	return recorder
}

func (r *Requester) Post(data []byte) *httptest.ResponseRecorder {
	router := r.server()

	recorder := httptest.NewRecorder()
	request, err := http.NewRequest("POST", r.URL, bytes.NewReader(data))

	Expect(err).To(BeNil())

	router.ServeHTTP(recorder, request)

	return recorder
}

func (r *Requester) Delete(params Params) *httptest.ResponseRecorder {
	router := r.server()

	url := r.parseParams(r.URL, params)
	recorder := httptest.NewRecorder()
	request, err := http.NewRequest("DELETE", url, nil)

	Expect(err).To(BeNil())

	router.ServeHTTP(recorder, request)

	return recorder
}

func (r *Requester) server() *gin.Engine {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.GET(r.URL, r.Handler)
	router.POST(r.URL, r.Handler)
	router.DELETE(r.URL, r.Handler)
	return router
}

func (r *Requester) parseParams(url string, params Params) string {
	for key, value := range params {
		url = strings.Replace(url, ":"+key, value, 1)
	}
	return url
}
func (r *Requester) parseQuery(url string, query Params) string {
	for key, value := range query {
		url += fmt.Sprintf("&%s=%s", key, value)
	}
	return url
}
