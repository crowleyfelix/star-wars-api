//+build integration

package controllers

import (
	"net/http"
	"net/http/httptest"
	"strconv"

	"github.com/crowleyfelix/star-wars-api/server/controllers/fixtures"
	"github.com/crowleyfelix/star-wars-api/server/models"
	. "github.com/crowleyfelix/star-wars-api/server/testing"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Planet CRUD", func() {

	var (
		id       int
		recorder *httptest.ResponseRecorder
	)

	Describe("When creating a planet", func() {

		var (
			requester = Requester{
				URL:     "/planets",
				Handler: Planet,
			}
			planet models.Planet
			actual fixtures.PlanetResponse
		)

		Context("with a valid planet", func() {
			BeforeEach(func() {
				LoadJSON("fixtures/planet.json", &planet)
				recorder = requester.Post(ToJSONBytes(planet))
			})

			It("should return created response", func() {
				Expect(recorder.Code).To(Equal(http.StatusCreated))

				LoadJSONFromBytes(recorder.Body.Bytes(), &actual)
				id = actual.Data.ID
				planet.ID = id
				Expect(actual.Data).To(Equal(planet))
			})
		})
		Context("with a invalid planet", func() {
			BeforeEach(func() {
				recorder = requester.Post(File("fixtures/planet-invalid.json"))
			})

			It("should return created response", func() {
				Expect(recorder.Code).To(Equal(http.StatusUnprocessableEntity))
			})
		})

	})

	Describe("When getting a planet", func() {

		var (
			requester = Requester{
				URL:     "/planets/:id",
				Handler: Planet,
			}
			expected fixtures.PlanetResponse
			response fixtures.PlanetResponse
		)

		BeforeEach(func() {
			LoadJSON("fixtures/planet-response.json", &expected)
			expected.Data.ID = id

			recorder = requester.Get(Params{"id": strconv.Itoa(id)}, nil)
		})

		It("should return planet", func() {
			Expect(recorder.Code).To(Equal(http.StatusOK))

			LoadJSONFromBytes(recorder.Body.Bytes(), &response)
			Expect(response).To(Equal(expected))
		})
	})

	Describe("When searching for planets", func() {

		var (
			requester = Requester{
				URL:     "/planets",
				Handler: Planets,
			}
			expected fixtures.PlanetsResponse
			response fixtures.PlanetsResponse
		)

		BeforeEach(func() {
			LoadJSON("fixtures/planets-response.json", &expected)
			expected.Data[0].ID = id
		})

		Context("with no parameter", func() {
			BeforeEach(func() {
				recorder = requester.Get(nil, nil)
			})

			It("should return a planet", func() {
				Expect(recorder.Code).To(Equal(http.StatusOK))

				LoadJSONFromBytes(recorder.Body.Bytes(), &response)
				Expect(response).To(Equal(expected))
			})
		})

		Context("with parameter name", func() {
			BeforeEach(func() {
				recorder = requester.Get(nil, Params{
					"name": expected.Data[0].Name,
				})
			})

			It("should return a planet", func() {
				Expect(recorder.Code).To(Equal(http.StatusOK))

				LoadJSONFromBytes(recorder.Body.Bytes(), &response)
				Expect(response.Data).To(Equal(expected.Data))
			})
		})
	})

	Describe("When deleting a planet", func() {

		var (
			requester = Requester{
				URL:     "/planets/:id",
				Handler: Planet,
			}
		)

		BeforeEach(func() {
			recorder = requester.Delete(Params{"id": strconv.Itoa(id)})
		})

		It("should return planets page", func() {
			Expect(recorder.Code).To(Equal(http.StatusOK))
		})
	})
})
