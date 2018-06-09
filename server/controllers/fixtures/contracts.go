package fixtures

import "github.com/crowleyfelix/star-wars-api/server/models"

type Response struct {
	*Page
	Data     interface{} `json:"data"`
	Messages []string    `json:"messages"`
}

type Page struct {
	Previous *string `json:"previous,omitempty"`
	Current  string  `json:"current,omitempty"`
	Next     *string `json:"next,omitempty"`
	Size     int     `json:"size,omitempty"`
	MaxSize  int     `json:"maxSize,omitempty"`
}

type PlanetResponse struct {
	Response
	Data models.Planet `json:"data"`
}

type PlanetsResponse struct {
	Response
	Data []models.Planet `json:"data"`
}
