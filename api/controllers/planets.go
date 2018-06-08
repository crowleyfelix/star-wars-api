package controllers

import (
	"github.com/crowleyfelix/star-wars-api/api/services"
	"github.com/gin-gonic/gin"
	"github.com/stone-payments/CaduGO/errors"
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
		p.fail(errors.NewInternalServer(err.Error()))
		return
	}

	p.ok(page)
}
