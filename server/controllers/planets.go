package controllers

import (
	"strconv"

	"github.com/crowleyfelix/star-wars-api/server/errors"
	"github.com/crowleyfelix/star-wars-api/server/models"
	"github.com/crowleyfelix/star-wars-api/server/services"
	"github.com/gin-gonic/gin"
)

type planets struct {
	baseController
	service services.Planet
}

//Planets handles with planet request
func Planets(c *gin.Context) {
	invokeMethod(c, &planets{
		baseController{c},
		services.NewPlanet(),
	})
}

func (p *planets) Get() {
	query := new(PlanetSearchParams)
	pagination := new(Pagination)

	if err := p.context.BindQuery(query); err != nil {
		p.fail(errors.NewBadRequest(err.Error()))
		return
	}

	if err := p.context.BindQuery(pagination); err != nil {
		p.fail(errors.NewBadRequest(err.Error()))
		return
	}

	page, err := p.service.Search(query.To(), pagination.To())

	if err != nil {
		p.fail(err)
		return
	}

	p.ok(page.Planets, page.Page)
}

type planet struct {
	planets
}

//Planet handles with planet request
func Planet(c *gin.Context) {
	invokeMethod(c, &planet{
		planets{
			baseController{c},
			services.NewPlanet(),
		},
	})
}

func (p *planet) Get() {
	id, err := p.intParam("id")
	if err != nil {
		p.fail(err)
		return
	}

	data, err := p.service.Get(*id)

	if err != nil {
		p.fail(err)
		return
	}

	p.ok(data, nil)
}

func (p *planet) Post() {

	var (
		data = new(models.Planet)
		err  errors.Error
	)

	if err := p.context.BindJSON(data); err != nil {
		p.fail(errors.NewUnprocessableEntity(err.Error()))
		return
	}

	validate := p.context.DefaultQuery("validate", "true")
	if validate, _ := strconv.ParseBool(validate); validate {
		if err := p.service.Validate(data); err != nil {
			p.fail(err)
			return
		}
	}

	if data, err = p.service.Create(data); err != nil {
		p.fail(err)
		return
	}

	p.created(data, nil)
}

func (p *planet) Delete() {
	id, err := p.intParam("id")
	if err != nil {
		p.fail(err)
		return
	}

	if err = p.service.Remove(*id); err != nil {
		p.fail(err)
		return
	}

	p.ok(nil, nil)
}
