package controllers

import (
	"net/http"

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
		c.JSON(http.StatusMethodNotAllowed, Response{"Method not allowed"})
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
	b.context.JSON(http.StatusMethodNotAllowed, Response{"Method not allowed!"})
}
