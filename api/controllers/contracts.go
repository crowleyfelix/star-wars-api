package controllers

import "github.com/crowleyfelix/star-wars-api/api/services"

type Response struct {
	Data     interface{} `json:"data"`
	Messages []string    `json:"messages"`
}

type Dependencie struct {
	Name  string  `json:"name"`
	Error *string `json:"error"`
}

type Pagination struct {
	Page int `form:"page"`
	Size int `form:"page_size"`
}

func (p *Pagination) To() *services.Pagination {
	pagination := new(services.Pagination)
	pagination.Page = p.Page
	pagination.Size = p.Size
	return pagination
}

type PlanetSearchParams struct {
	ID      *int    `form:"id"`
	Name    *string `form:"name"`
	Climate *string `form:"climate"`
	Terrain *string `form:"terrain"`
}

func (q *PlanetSearchParams) To() *services.PlanetSearchParams {
	query := new(services.PlanetSearchParams)
	query.ID = q.ID
	query.Name = q.Name
	query.Climate = q.Climate
	query.Terrain = q.Terrain
	return query
}
