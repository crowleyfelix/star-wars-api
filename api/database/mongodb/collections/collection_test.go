package collections

import (
	"github.com/crowleyfelix/star-wars-api/api/errors"

	"github.com/bouk/monkey"
	"github.com/crowleyfelix/star-wars-api/api/database/mongodb"
	"github.com/crowleyfelix/star-wars-api/api/database/mongodb/mocks"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"gopkg.in/mgo.v2"
)

var _ = Describe("Collection", func() {

	var (
		coll     *collection
		mockPool = new(mocks.SessionPool)
	)

	BeforeEach(func() {
		coll = new(collection)
		monkey.Patch(mongodb.Pool, func() mongodb.SessionPool { return mockPool })
	})
	AfterEach(func() { monkey.UnpatchAll() })

	Describe("execute(): When executing operation", func() {
		var (
			actualErr     error
			mockSession   = new(mgo.Session)
			mockOperation = func(c *mgo.Collection) error {
				return errors.Build(0)
			}
		)

		JustBeforeEach(func() {
			actualErr = coll.execute(mockOperation)
		})
		Context("and failed creating session", func() {
			BeforeEach(func() {
				mockPool.On("Session").Return(nil, errors.Build(0)).Once()
			})
			It("should return an error", func() {
				Expect(actualErr).To(BeAssignableToTypeOf(new(errors.InternalServer)))
			})
		})
		Context("and success creating session", func() {
			BeforeEach(func() {
				mockPool.On("Session").Return(mockSession, nil).Once()
				mockPool.On("Release", mockSession).Return().Once()
			})
			It("should return operation result", func() {
				Expect(actualErr).To(BeAssignableToTypeOf(new(errors.InternalServer)))
			})
		})
	})
})
