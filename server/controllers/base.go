package controllers

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	"github.com/aphistic/gomol"

	"github.com/crowleyfelix/star-wars-api/server/models"

	"github.com/crowleyfelix/star-wars-api/server/errors"

	"github.com/gin-gonic/gin"
)

const (
	paginationQuery   = "%s?page=%d&page_size=%d%s"
	paginationPattern = "((page=\\d+&?)|(page_size=\\d+&?))"
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
	gomol.Debug("Sending fail response")
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
	gomol.Debug("Sending ok response")

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
	gomol.Debug("Sending created response")

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
		gomol.Warning("Page is nil, nothing to calculate")
		return nil
	}

	url := b.context.(*gin.Context).Request.URL
	compiled, _ := regexp.Compile(paginationPattern)
	query := compiled.ReplaceAllString(url.RawQuery, "")

	pg := &Page{
		Current: fmt.Sprintf(paginationQuery, url.Path, page.Current, page.MaxSize, query),
		MaxSize: page.MaxSize,
		Size:    page.Size,
	}

	if page.Previous != nil {
		previous := fmt.Sprintf(paginationQuery, url.Path, *page.Previous, page.MaxSize, query)
		pg.Previous = &previous
	}

	if page.Next != nil {
		next := fmt.Sprintf(paginationQuery, url.Path, *page.Next, page.MaxSize, query)
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
