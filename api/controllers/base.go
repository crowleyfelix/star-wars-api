package controllers

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	"github.com/crowleyfelix/star-wars-api/api/models"

	"github.com/crowleyfelix/star-wars-api/api/errors"

	"github.com/gin-gonic/gin"
)

const (
	paginationQuery = "%s&page=%d&page_size=%d"
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

func (b *baseController) ok(data interface{}, page *models.Page) {

	if data == nil {
		data = make(map[string]interface{})
	}

	b.context.JSON(http.StatusOK, Response{
		Page:     b.calculatePage(page),
		Data:     data,
		Messages: make([]string, 0),
	})
}

func (b *baseController) created(data interface{}, page *models.Page) {

	if data == nil {
		data = make(map[string]interface{})
	}

	b.context.JSON(http.StatusCreated, Response{
		Data:     data,
		Messages: make([]string, 0),
	})
}

func (b *baseController) calculatePage(page *models.Page) *Page {

	if page == nil {
		return nil
	}

	url := b.context.(*gin.Context).Request.URL

	compiled, _ := regexp.Compile("((page=\\d+&?)|(page_size=\\d+&?))")
	query := compiled.ReplaceAllString(url.RawQuery, "")

	path := fmt.Sprintf("%s?%s", url.Path, query)

	pg := &Page{
		Current: fmt.Sprintf(paginationQuery, path, page.Current, page.MaxSize),
		MaxSize: page.MaxSize,
		Size:    page.Size,
	}

	if page.Previous != nil {
		previous := fmt.Sprintf(paginationQuery, path, page.Previous, page.MaxSize)
		pg.Previous = &previous
	}

	if page.Next != nil {
		next := fmt.Sprintf(paginationQuery, path, page.Next, page.MaxSize)
		pg.Next = &next
	}

	return pg
}

func (b *baseController) intParam(name string) (*int, errors.Error) {
	n, err := strconv.Atoi(b.context.Param(name))

	if err != nil {
		return nil, errors.NewBadRequest(fmt.Sprintf("invalid param sent %s", name))
	}

	return &n, nil
}
