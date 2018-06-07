//+build integration

package mongodb

import (
	"strings"

	"github.com/crowleyfelix/star-wars-api/src/mongodb/models"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/satori/go.uuid"
)

var _ = Describe("PlanetCollection", func() {

	var (
		coll   = NewPlanetCollection()
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

	Describe("Find(): When finding planet", func() {
		var (
			results []models.Planet
			err     error
		)

		JustBeforeEach(func() {
			query := &PlanetSearchQuery{
				Name: &planet.Name,
			}
			results, err = coll.Find(query)
		})
		AfterEach(func() {
			if len(results) > 0 {
				id = results[0].ID
			}
		})
		It("should get planet", func() {
			Expect(err).To(BeNil())
			Expect(results).To(HaveLen(1))
		})
	})

	Describe("FindByID(): When finding planet by id", func() {
		var (
			actual *models.Planet
			err    error
		)

		JustBeforeEach(func() {
			actual, err = coll.FindByID(planet.ID)
		})
		It("should get planet", func() {
			Expect(err).To(BeNil())
			Expect(actual).To(BeEquivalentTo(planet))
		})
	})

	Describe("List(): When listing planets", func() {
		var (
			pagination = &Pagination{
				Page: 1,
				Size: 10,
			}
		)

		var (
			page *models.PlanetPage
			err  error
		)

		JustBeforeEach(func() {
			page, err = coll.List(pagination)
		})
		It("should list planets paged", func() {
			Expect(err).To(BeNil())
			Expect(len(page.Planets) < 11).To(BeTrue())
			Expect(page.MaxSize).To(Equal(pagination.Size))
			Expect(page.Current).To(Equal(pagination.Page))
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
				Expect(err).ToNot(BeNil())
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
				Expect(err).ToNot(BeNil())
			})
		})
	})
})
