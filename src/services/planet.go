package services

import (
	"sync"

	"github.com/crowleyfelix/star-wars-api/src/clients/swapi"
	"github.com/crowleyfelix/star-wars-api/src/models"
	"github.com/crowleyfelix/star-wars-api/src/mongodb"
)

//Planet exposes necessary methods of a planet service
type Planet interface {
	Create(*models.Planet) error
	Get(int) (*models.Planet, error)
	Search(*PlanetSearchParams, *Pagination) (*models.PlanetPage, error)
	Remove(int) error
}

type planet struct {
	Planet
	client   swapi.Client
	database mongodb.PlanetCollection
}

//NewPlanet returns a new planet service
func NewPlanet() Planet {
	return &planet{
		client:   swapi.New(),
		database: mongodb.NewPlanetCollection(),
	}
}

func (p *planet) Create(data *models.Planet) error {
	return p.database.Insert(data.To())
}

func (p *planet) Get(id int) (*models.Planet, error) {
	raw, err := p.database.FindByID(id)

	if err != nil {
		return nil, err
	}

	var data models.Planet
	data.From(raw)

	return &data, p.fetchFilms([]models.Planet{data})
}

func (p *planet) Search(params *PlanetSearchParams, pagination *Pagination) (*models.PlanetPage, error) {
	raw, err := p.database.Find(&params.PlanetSearchQuery, &pagination.Pagination)

	if err != nil {
		return nil, err
	}

	page := new(models.PlanetPage)
	page.From(raw)
	return page, p.fetchFilms(page.Planets)
}

func (p *planet) Remove(id int) error {
	return p.database.Delete(id)
}

func (p *planet) fetchFilms(planets []models.Planet) error {
	group := new(sync.WaitGroup)
	mutex := new(sync.Mutex)

	var err error

	for i := range planets {
		group.Add(1)

		planet := planets[i]
		go func() {
			defer mutex.Unlock()

			films, e := p.client.PlanetFilms(planet.Name)

			if e != nil {
				err = e
			}

			mutex.Lock()
			planet.Films.From(films)

			group.Done()
		}()
	}

	group.Wait()
	return err
}
