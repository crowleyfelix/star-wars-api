package controllers

import (
	"net/http"
	"reflect"

	"github.com/bouk/monkey"

	"github.com/gin-gonic/gin"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/crowleyfelix/star-wars-api/api/controllers/mocks"
)

var _ = Describe("Base", func() {

	var (
		context *gin.Context
		handler = new(mocks.Controller)
	)

	BeforeEach(func() {
		context = &gin.Context{
			Request: new(http.Request),
		}
	})
	AfterEach(func() {
		monkey.UnpatchAll()
	})

	Describe("invokeMethod", func() {
		Context("When received a GET http method", func() {
			BeforeEach(func() {
				context.Request.Method = http.MethodGet
				handler.On("Get").Return().Once()
			})
			It("should call handler Get method", func() {
				invokeMethod(context, handler)
			})
		})
		Context("When received a POST http method", func() {
			BeforeEach(func() {
				context.Request.Method = http.MethodPost
				handler.On("Post").Return().Once()
			})
			It("should call handler Post method", func() {
				invokeMethod(context, handler)
			})
		})
		Context("When received a PUT http method", func() {
			BeforeEach(func() {
				context.Request.Method = http.MethodPut
				handler.On("Put").Return().Once()
			})
			It("should call handler Put method", func() {
				invokeMethod(context, handler)
			})
		})
		Context("When received a DELETE http method", func() {
			BeforeEach(func() {
				context.Request.Method = http.MethodDelete
				handler.On("Delete").Return().Once()
			})
			It("should call handler Delete method", func() {
				invokeMethod(context, handler)
			})
		})
		Context("When received an invalid http method", func() {
			BeforeEach(func() {
				context.Request.Method = http.MethodTrace
				monkey.PatchInstanceMethod(reflect.TypeOf(context), "JSON", func(_ *gin.Context, code int, obj interface{}) {
					Expect(code).To(Equal(http.StatusMethodNotAllowed))
				})
			})
			It("should not call any handler", func() {
				invokeMethod(context, handler)
			})
		})
	})
})

var _ = Describe("Controller", func() {

	var (
		context = new(mocks.RequestContext)
		base    *baseController
	)

	BeforeEach(func() {
		base = &baseController{context}

		context.On("JSON", http.StatusMethodNotAllowed, Response{Message: "Method not allowed!"}).
			Return().Once()
	})
	Describe("Get()", func() {
		It("should respond a method not allowed status", func() {
			base.Get()
		})
	})
	Describe("Post()", func() {
		It("should respond a method not allowed status", func() {
			base.Post()
		})
	})
	Describe("Put()", func() {
		It("should respond a method not allowed status", func() {
			base.Put()
		})
	})
	Describe("Delete()", func() {
		It("should respond a method not allowed status", func() {
			base.Delete()
		})
	})
})
