//+build !integration

package mongodb

import (
	"errors"

	"github.com/crowleyfelix/star-wars-api/src/mongodb/mocks"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"gopkg.in/mgo.v2"
)

var _ = Describe("Collection", func() {

	var (
		coll     *collection
		mockPool = new(mocks.SessionManager)
	)

	BeforeEach(func() {
		coll = new(collection)
		Pool = mockPool
	})

	Describe("execute", func() {
		var (
			actualErr     error
			expectedErr   error
			mockSession   = new(mgo.Session)
			mockOperation = func(c *mgo.Collection) error {
				return expectedErr
			}
		)

		JustBeforeEach(func() {
			actualErr = coll.execute(mockOperation)
		})
		Context("When executing operation", func() {
			Context("and failed creating session", func() {
				BeforeEach(func() {
					expectedErr = errors.New("error")
					mockPool.On("Session").Return(nil, expectedErr).Once()
				})
				It("should return an error", func() {
					Expect(actualErr).To(Equal(expectedErr))
				})
			})
			Context("and success creating session", func() {
				BeforeEach(func() {
					mockPool.On("Session").Return(mockSession, nil).Once()
					mockPool.On("Release", mockSession).Return().Once()
				})
				It("should return operation result", func() {
					Expect(actualErr).To(Equal(expectedErr))
				})
			})
		})
	})
})
