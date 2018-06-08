package services

import (
	swapi "github.com/crowleyfelix/star-wars-api/api/clients/swapi"
	swapiMocks "github.com/crowleyfelix/star-wars-api/api/clients/swapi/mocks"
	mongodbMocks "github.com/crowleyfelix/star-wars-api/api/database/mongodb/collections/mocks"
	mongoModels "github.com/crowleyfelix/star-wars-api/api/database/mongodb/models"
	"github.com/crowleyfelix/star-wars-api/api/errors"
	"github.com/crowleyfelix/star-wars-api/api/models"

	. "github.com/crowleyfelix/star-wars-api/api/testing"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Planet", func() {

	var (
		pl           planet
		databaseMock *mongodbMocks.Planets
		clientMock   *swapiMocks.Client
	)

	BeforeEach(func() {
		databaseMock = new(mongodbMocks.Planets)
		clientMock = new(swapiMocks.Client)

		pl = planet{
			database: databaseMock,
			client:   clientMock,
		}
	})

	Describe("NewPlanet()", func() {

		BeforeEach(func() {
			pl = *NewPlanet().(*planet)
		})

		It("should set dependencies", func() {
			Expect(pl.client).ToNot(BeNil())
			Expect(pl.database).ToNot(BeNil())
		})
	})

	Describe("Create(): When creating a planet", func() {
		var (
			modelPlanet models.Planet
			mongoPlanet mongoModels.Planet
		)

		var (
			actualError   errors.Error
			expectedError error
		)

		BeforeEach(func() {
			LoadJSON("fixtures/model-planet.json", &modelPlanet)
			LoadJSON("fixtures/mongo-planet.json", &mongoPlanet)

			expectedError = errors.Build(0)
			databaseMock.On("Insert", &mongoPlanet).Return(expectedError).Once()

			actualError = pl.Create(&modelPlanet)
		})

		It("should insert into database", func() {
			Expect(actualError).To(Equal(expectedError))
		})
	})

	Describe("Get(): When getting a planet", func() {
		var (
			id = 1
		)

		var (
			actualPlanet   *models.Planet
			actualError    errors.Error
			expectedPlanet *models.Planet
			expectedError  errors.Error
		)

		JustBeforeEach(func() {
			actualPlanet, actualError = pl.Get(id)
		})

		Context("and failed on find planet in database", func() {
			BeforeEach(func() {
				expectedError = errors.Build(0)
				databaseMock.On("FindByID", id).Return(nil, expectedError).Once()
			})

			It("should return an error", func() {
				Expect(actualPlanet).To(BeNil())
				Expect(actualError).To(Equal(expectedError))
			})
		})
		Context("and sucess on getting planet from database", func() {
			var (
				mongoPlanet mongoModels.Planet
				swapiFilms  []swapi.Film
			)

			BeforeEach(func() {
				LoadJSON("fixtures/mongo-planet.json", &mongoPlanet)
				databaseMock.On("FindByID", id).Return(&mongoPlanet, nil).Once()
			})
			Context("and failed on fetch planet films", func() {
				BeforeEach(func() {
					expectedError = errors.Build(0)
					clientMock.On("PlanetFilms", mongoPlanet.Name).Return(nil, expectedError).Once()
				})
				It("should return an error", func() {
					Expect(actualPlanet).To(BeNil())
					Expect(actualError).To(Equal(expectedError))
				})
			})
			Context("and success on fetch planet films", func() {
				BeforeEach(func() {
					LoadJSON("fixtures/model-planet.json", &expectedPlanet)
					LoadJSON("fixtures/swapi-films.json", &swapiFilms)

					clientMock.On("PlanetFilms", mongoPlanet.Name).Return(swapiFilms, nil).Once()
				})
				It("should return an error", func() {
					Expect(actualPlanet).To(Equal(expectedPlanet))
					Expect(actualError).To(BeNil())
				})
			})
		})
	})

	Describe("Search(): When searching for a planet", func() {
		var (
			params     = new(PlanetSearchParams)
			pagination = new(Pagination)
		)

		var (
			actualPage    *models.PlanetPage
			actualError   errors.Error
			expectedPage  *models.PlanetPage
			expectedError error
		)

		JustBeforeEach(func() {
			actualPage, actualError = pl.Search(params, pagination)
		})

		Context("and failed on find planets in database", func() {
			BeforeEach(func() {
				expectedError = errors.Build(0)
				databaseMock.
					On("Find", &params.PlanetSearchQuery, &pagination.Pagination).
					Return(nil, expectedError).Once()
			})

			It("should return an error", func() {
				Expect(actualPage).To(BeNil())
				Expect(actualError).To(Equal(expectedError))
			})
		})
		Context("and sucess on getting planets from database", func() {
			var (
				mongoPage  mongoModels.PlanetPage
				swapiFilms []swapi.Film
			)

			BeforeEach(func() {
				LoadJSON("fixtures/mongo-planet-page.json", &mongoPage)
				databaseMock.
					On("Find", &params.PlanetSearchQuery, &pagination.Pagination).
					Return(&mongoPage, nil).Once()
			})
			Context("and failed on fetch planet films", func() {
				BeforeEach(func() {
					expectedError = errors.Build(0)
					clientMock.
						On("PlanetFilms", mongoPage.Planets[0].Name).
						Return(nil, expectedError).Once()
				})
				It("should return an error", func() {
					Expect(actualPage).To(BeNil())
					Expect(actualError).To(Equal(expectedError))
				})
			})
			Context("and success on fetch planet films", func() {
				BeforeEach(func() {
					LoadJSON("fixtures/model-planet-page.json", &expectedPage)
					LoadJSON("fixtures/swapi-films.json", &swapiFilms)

					clientMock.
						On("PlanetFilms", mongoPage.Planets[0].Name).
						Return(swapiFilms, nil).Once()
				})
				It("should return an error", func() {
					Expect(actualPage).To(Equal(expectedPage))
					Expect(actualError).To(BeNil())
				})
			})
		})
	})

	Describe("Remove(): When removing a planet", func() {
		var (
			id = 12
		)

		var (
			actualError   errors.Error
			expectedError error
		)

		BeforeEach(func() {
			expectedError = errors.Build(0)
			databaseMock.On("Delete", id).Return(expectedError).Once()

			actualError = pl.Remove(id)
		})

		It("should remove from database", func() {
			Expect(actualError).To(Equal(expectedError))
		})
	})
})
