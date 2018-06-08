package controllers

import (
	"net/http"

	"github.com/crowleyfelix/star-wars-api/api/errors"

	"github.com/gin-gonic/gin"
)

//Controller exposes api controller methdos
type Controller interface {
	Get()
	Post()
	Put()
	Delete()
}

func invokeMethod(c *gin.Context, h Controller) {
	switch c.Request.Method {
	case http.MethodGet:
		h.Get()
	case http.MethodPost:
		h.Post()
	case http.MethodPut:
		h.Put()
	case http.MethodDelete:
		h.Delete()
	default:
		c.JSON(http.StatusMethodNotAllowed, Response{Messages: []string{"Method not allowed"}})
	}
}

type baseController struct {
	context RequestContext
}

func (b *baseController) Get() {
	b.notAllowed()
}

func (b *baseController) Post() {
	b.notAllowed()
}

func (b *baseController) Put() {
	b.notAllowed()
}

func (b *baseController) Delete() {
	b.notAllowed()
}

func (b *baseController) notAllowed() {
	b.fail(errors.NewMethodNotAllowed())
}

func (b *baseController) fail(err errors.Error) {
	statusCode := http.StatusInternalServerError

	if value, ok := err.(errors.HTTPError); ok {
		statusCode = value.StatusCode()
	}

	b.context.JSON(statusCode, Response{
		Data:     make(map[string]interface{}),
		Messages: err.Messages(),
	})
}

func (b *baseController) ok(data interface{}) {
	b.context.JSON(http.StatusOK, Response{
		Data:     data,
		Messages: make([]string, 0),
	})
}
