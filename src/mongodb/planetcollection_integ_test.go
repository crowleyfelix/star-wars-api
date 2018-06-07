//+build integration

package mongodb

import (
	"strings"

	"github.com/crowleyfelix/star-wars-api/src/mongodb/models"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("PlanetCollection", func() {

	var (
		coll   = NewPlanetCollection()
		planet *models.Planet
	)

	BeforeEach(func() {
		planet = &models.Planet{
			Name:    "Alderaan",
			Terrain: "grasslands, mountains",
			Climate: "temperate",
		}
	})

	Describe("When inserting a planet", func() {
		var err error

		JustBeforeEach(func() {
			err = coll.Insert(planet)
		})
		Context("with unexistent data", func() {
			It("should set id", func() {
				Expect(err).To(BeNil())
			})
		})
		Context("with duplicated name", func() {
			BeforeEach(func() {
				planet.Name = strings.ToLower(planet.Name)
			})
			It("should return an error", func() {
				Expect(err).ToNot(BeNil())
			})
		})
	})

	Describe("When finding planet", func() {
		var (
			actual *models.Planet
			err    error
		)

		JustBeforeEach(func() {
			actual, err = coll.Find(0)
		})
		It("should get planet", func() {
			Expect(err).To(BeNil())
			Expect(actual).To(BeEquivalentTo(planet))
		})
	})

	Describe("When listing planets", func() {
		var (
			pagination = &Pagination{
				Page: 1,
				Size: 10,
			}
		)

		var (
			actual []models.Planet
			err    error
		)

		JustBeforeEach(func() {
			actual, err = coll.List(pagination)
		})
		It("should list planets paged", func() {
			Expect(err).To(BeNil())
			Expect(actual).To(HaveLen(1))
			Expect(actual[0]).To(BeEquivalentTo(*planet))
		})
	})

	Describe("When updating a planet", func() {
		var (
			err error
		)

		JustBeforeEach(func() {
			planet.Name = "Hoth"
			err = coll.Update(planet)
		})
		It("should update planet", func() {
			Expect(err).To(BeNil())
		})
	})

	Describe("When deleting a planet", func() {
		var (
			err error
		)

		JustBeforeEach(func() {
			err = coll.Delete(planet.ID)
		})
		It("should not return an error", func() {
			Expect(err).To(BeNil())
		})
	})
})
