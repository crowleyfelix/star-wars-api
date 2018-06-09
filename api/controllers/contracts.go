package controllers

import (
	"strconv"

	"github.com/crowleyfelix/star-wars-api/api/services"
)

const (
	defaultPage     = 1
	defaultPageSize = 10
)

type Response struct {
	*Page
	Data     interface{} `json:"data"`
	Messages []string    `json:"messages"`
}

type Dependence struct {
	Name  string  `json:"name"`
	Error *string `json:"error"`
}

type Page struct {
	Previous *string `json:"previous,omitempty"`
	Current  string  `json:"current,omitempty"`
	Next     *string `json:"next,omitempty"`
	Size     int     `json:"size,omitempty"`
	MaxSize  int     `json:"maxSize,omitempty"`
}

type Pagination struct {
	Page int `form:"page"`
	Size int `form:"page_size"`
}

func (p *Pagination) To() *services.Pagination {
	pagination := new(services.Pagination)

	if p.Page == 0 {
		pagination.Page = defaultPage
	} else {
		pagination.Page = p.Page
	}

	if p.Size == 0 {
		pagination.Size = defaultPageSize
	} else {
		pagination.Size = p.Size
	}

	return pagination
}

type PlanetSearchParams struct {
	ID      string `form:"id"`
	Name    string `form:"name"`
	Climate string `form:"climate"`
	Terrain string `form:"terrain"`
}

func (q *PlanetSearchParams) To() *services.PlanetSearchParams {
	query := new(services.PlanetSearchParams)

	if q.ID != "" {
		id, err := strconv.Atoi(q.ID)

		if err != nil {
			query.ID = &id
		}
	}

	if q.Name != "" {
		query.Name = &q.Name
	}

	if q.Climate != "" {
		query.Climate = &q.Climate
	}

	if q.Terrain != "" {
		query.Terrain = &q.Terrain
	}

	return query
}
