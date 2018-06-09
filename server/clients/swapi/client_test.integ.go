//+build integration

package swapi

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Client", func() {

	var cl Client

	BeforeEach(func() {
		cl = New()
	})

	Describe("PlanetFilms(): When getting planet films", func() {
		var (
			films []Film
			err   error
		)

		BeforeEach(func() {
			films, err = cl.PlanetFilms("Alderaan")
		})

		It("should return all of them", func() {
			Expect(films).ToNot(BeEmpty())
		})
	})
})
