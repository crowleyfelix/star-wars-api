//+build integration

package collections

import (
	"strings"

	"github.com/crowleyfelix/star-wars-api/server/database/mongodb/models"
	"github.com/crowleyfelix/star-wars-api/server/errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/satori/go.uuid"
)

var _ = Describe("Planets", func() {

	var (
		coll   = NewPlanets()
		planet *models.Planet
		id     int
		name   string
	)

	BeforeSuite(func() {
		uid, err := uuid.NewV4()
		Expect(err).To(BeNil())

		name = uid.String()
	})

	BeforeEach(func() {
		planet = &models.Planet{
			ID:      id,
			Name:    name,
			Terrain: "grasslands, mountains",
			Climate: "temperate",
		}
	})

	Describe("Insert(): When inserting a planet", func() {
		var (
			actual *models.Planet
			err    error
		)

		JustBeforeEach(func() {
			actual, err = coll.Insert(planet)

			if actual != nil {
				id = actual.ID
			}
		})
		Context("with unexistent data", func() {
			It("should set id", func() {
				Expect(err).To(BeNil())
				Expect(actual).To(Equal(planet))
			})
		})
		Context("with duplicated name", func() {
			BeforeEach(func() {
				planet.Name = strings.ToLower(planet.Name)
			})
			It("should return an error", func() {
				Expect(err).To(BeAssignableToTypeOf(new(errors.UnprocessableEntity)))
			})
		})
	})

	Describe("Find(): When finding planets", func() {
		var (
			pagination *Pagination
		)

		var (
			actual *models.PlanetPage
			expect *models.PlanetPage
			err    error
		)

		BeforeEach(func() {
			pagination = &Pagination{
				Page: 1,
				Size: 1,
			}
			query := &PlanetSearchQuery{
				Name: &planet.Name,
			}

			actual, err = coll.Find(query, pagination)

			expect = &models.PlanetPage{
				Page: &models.Page{
					MaxSize:  1,
					Size:     1,
					Current:  1,
					Previous: nil,
					Next:     nil,
				},
				Planets: []models.Planet{*planet},
			}
		})
		It("should get planet", func() {
			Expect(err).To(BeNil())
			Expect(actual).To(Equal(expect))
		})
	})

	Describe("FindByID(): When finding planet by id", func() {
		var (
			actual *models.Planet
			err    error
		)

		BeforeEach(func() {
			actual, err = coll.FindByID(planet.ID)
		})
		It("should get planet", func() {
			Expect(err).To(BeNil())
			Expect(actual).To(BeEquivalentTo(planet))
		})
	})

	Describe("Update(): When updating a planet", func() {
		var (
			changed *models.Planet
			err     error
		)

		JustBeforeEach(func() {
			err = coll.Update(planet)
			changed, _ = coll.FindByID(planet.ID)
		})
		Context("with existent entity", func() {
			BeforeEach(func() {
				uid, err := uuid.NewV4()
				Expect(err).To(BeNil())
				planet.Name = uid.String()
			})
			It("should update planet", func() {
				Expect(err).To(BeNil())
				Expect(changed).To(Equal(planet))
			})
		})
		Context("with unexistent entity", func() {
			BeforeEach(func() {
				planet.ID = planet.ID + 1
			})
			It("should return an error", func() {
				Expect(err).To(BeAssignableToTypeOf(new(errors.NotFound)))
			})
		})
	})

	Describe("Delete(): When deleting a planet", func() {
		var (
			err error
		)

		JustBeforeEach(func() {
			err = coll.Delete(planet.ID)
		})
		Context("with existent entity", func() {
			It("should not return an error", func() {
				Expect(err).To(BeNil())
			})
		})
		Context("with unexistent entity", func() {
			BeforeEach(func() {
				planet.ID = planet.ID + 1
			})
			It("should return an error", func() {
				Expect(err).To(BeAssignableToTypeOf(new(errors.NotFound)))
			})
		})
	})
})
