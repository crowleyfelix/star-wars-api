//+build integration

package controllers

import (
	"net/http"
	"net/http/httptest"

	. "github.com/crowleyfelix/star-wars-api/api/testing"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Planets", func() {

	Describe("Get", func() {

		var (
			requester = Requester{
				URL:     "/planets",
				Handler: Planets,
			}
			recorder *httptest.ResponseRecorder
		)

		BeforeEach(func() {
			recorder = requester.Get()
		})

		It("should return planets page", func() {
			Expect(recorder.Code).To(Equal(http.StatusOK))
		})
	})
})
