package services

import (
	"fmt"
	"sync"

	"github.com/aphistic/gomol"

	"github.com/crowleyfelix/star-wars-api/api/clients/swapi"
	mongodb "github.com/crowleyfelix/star-wars-api/api/database/mongodb/collections"
	"github.com/crowleyfelix/star-wars-api/api/errors"
	"github.com/crowleyfelix/star-wars-api/api/models"
)

//Planet exposes necessary methods of a planet service
type Planet interface {
	Create(*models.Planet) (*models.Planet, errors.Error)
	Get(int) (*models.Planet, errors.Error)
	Search(*PlanetSearchParams, *Pagination) (*models.PlanetPage, errors.Error)
	Remove(int) errors.Error
	Validate(*models.Planet) errors.Error
}

type planet struct {
	client   swapi.Client
	database mongodb.Planets
}

//NewPlanet returns a new planet service
func NewPlanet() Planet {
	return &planet{
		client:   swapi.New(),
		database: mongodb.NewPlanets(),
	}
}

func (p *planet) Create(data *models.Planet) (*models.Planet, errors.Error) {
	attrs := gomol.NewAttrsFromMap(map[string]interface{}{
		"planet": data,
	})

	gomol.Infom(attrs, "Inserting planet on database")
	result, err := p.database.Insert(data.To())

	if err != nil {
		return nil, err
	}

	data.From(result)
	return data, nil
}

func (p *planet) Get(id int) (*models.Planet, errors.Error) {
	gomol.Infof("Finding planet %d on database", id)
	raw, err := p.database.FindByID(id)

	if err != nil {
		return nil, err
	}

	var data models.Planet
	data.From(raw)
	data.Films = make(models.Films, 0)

	var planets []*models.Planet
	planets = append(planets, &data)

	if err = p.fetchFilms(planets); err != nil {
		return nil, err
	}

	return &data, nil
}

func (p *planet) Search(params *PlanetSearchParams, pagination *Pagination) (*models.PlanetPage, errors.Error) {
	attrs := gomol.NewAttrsFromMap(map[string]interface{}{
		"params":     params,
		"pagination": pagination,
	})

	gomol.Infom(attrs, "Searching planet on database")

	raw, err := p.database.Find(&params.PlanetSearchQuery, &pagination.Pagination)

	if err != nil {
		return nil, err
	}

	page := new(models.PlanetPage)
	page.From(raw)

	planets := make([]*models.Planet, 0)
	for i := range page.Planets {
		planets = append(planets, &page.Planets[i])
	}

	if err = p.fetchFilms(planets); err != nil {
		return nil, err
	}

	return page, nil
}

func (p *planet) Remove(id int) errors.Error {
	gomol.Infof("Deleting planet %d on database", id)
	return p.database.Delete(id)
}

func (p *planet) Validate(data *models.Planet) errors.Error {
	gomol.Info("Validating planet")
	_, err := p.client.Planet(data.Name)

	if _, ok := err.(*errors.NotFound); ok {
		return errors.NewUnprocessableEntity(fmt.Sprintf("The planet %s is invalid", data.Name))
	}

	return err
}

func (p *planet) fetchFilms(planets []*models.Planet) errors.Error {
	gomol.Debug("Fetching planets from swapi")

	group := new(sync.WaitGroup)
	mutex := new(sync.Mutex)

	var err errors.Error

	for i := range planets {
		group.Add(1)

		go func(planet *models.Planet) {
			defer mutex.Unlock()
			defer group.Done()

			attrs := gomol.NewAttrsFromMap(map[string]interface{}{
				"planet": planet,
			})
			gomol.Infom(attrs, "Fetching planet from swapi")

			films, e := p.client.PlanetFilms(planet.Name)

			if e != nil {
				err = e
			}

			mutex.Lock()
			planet.Films.From(films)
		}(planets[i])
	}

	group.Wait()
	return err
}
