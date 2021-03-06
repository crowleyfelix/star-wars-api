package models

import (
	"github.com/crowleyfelix/star-wars-api/server/database/mongodb/models"
)

//Planet represents a star wars planet
type Planet struct {
	ID      int    `json:"id"`
	Name    string `json:"name" binding:"required"`
	Climate string `json:"climate" binding:"required"`
	Terrain string `json:"terrain" binding:"required"`
	Films   Films  `json:"films"`
}

//From maps database planet model to application model
func (p *Planet) From(raw *models.Planet) {
	p.ID = raw.ID
	p.Name = raw.Name
	p.Climate = raw.Climate
	p.Terrain = raw.Terrain
}

//To maps application planet model to database model
func (p *Planet) To() *models.Planet {
	return &models.Planet{
		ID:      p.ID,
		Name:    p.Name,
		Climate: p.Climate,
		Terrain: p.Terrain,
	}
}

//PlanetPage represents a star wars planet page
type PlanetPage struct {
	*Page
	Planets []Planet `json:"planets"`
}

//From maps database planet model to application model
func (p *PlanetPage) From(raw *models.PlanetPage) {
	temp := Page(*raw.Page)
	p.Page = &temp
	p.Planets = make([]Planet, 0)
	for _, item := range raw.Planets {
		var planet Planet
		planet.From(&item)
		p.Planets = append(p.Planets, planet)
	}
}
