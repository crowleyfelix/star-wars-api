package collections

import (
	"fmt"

	"github.com/crowleyfelix/star-wars-api/server/errors"

	"github.com/aphistic/gomol"
	"github.com/crowleyfelix/star-wars-api/server/configuration"
	"github.com/crowleyfelix/star-wars-api/server/database/mongodb/models"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//Planets exposes methods of planet CRUD operations
type Planets interface {
	Insert(*models.Planet) (*models.Planet, errors.Error)
	Find(*PlanetSearchQuery, *Pagination) (*models.PlanetPage, errors.Error)
	FindByID(int) (*models.Planet, errors.Error)
	Update(*models.Planet) errors.Error
	Delete(int) errors.Error
}

type planets struct {
	collection
}

//NewPlanets returns new instance of planet collection
func NewPlanets() Planets {
	gomol.Debug("Creating new mongodb planet collection manager")

	config := configuration.Get().MongoDB
	return &planets{
		collection{
			DataBase:   config.Database,
			Collection: "planets",
			CounterID:  "planet_id",
		},
	}
}

func (pr *planets) Insert(planet *models.Planet) (*models.Planet, errors.Error) {
	gomol.Debug("Inserting planet on mongodb")

	err := pr.execute(func(col *mgo.Collection) error {
		var err errors.Error
		planet.ID, err = pr.calculateNextID(col.Database)

		if err != nil {
			gomol.Errorf("Failed on calculating next id: %s", err.Error())
			return err
		}

		return col.Insert(planet)
	})

	if err != nil {
		return nil, err
	}

	return planet, nil
}

func (pr *planets) FindByID(id int) (*models.Planet, errors.Error) {
	gomol.Debug("Finding planet by id on mongodb")

	query := &PlanetSearchQuery{
		ID: &id,
	}
	pagination := &Pagination{
		Page: 1,
		Size: 1,
	}

	page, err := pr.Find(query, pagination)

	if err != nil {
		gomol.Errorf("Failed on finding planet id %d on database: %s", id, err.Error())
		return nil, err
	}

	if page.Size == 0 {
		gomol.Errorf("Planet id %d was not found on database", id)
		return nil, errors.NewNotFound(fmt.Sprintf("Planet id %d was not found", id))
	}

	return &page.Planets[0], err
}

func (pr *planets) Find(query *PlanetSearchQuery, pagination *Pagination) (*models.PlanetPage, errors.Error) {
	gomol.Debug("Finding planet on mongodb")

	var (
		err  errors.Error
		page = new(models.PlanetPage)
	)

	err = pr.execute(func(col *mgo.Collection) error {
		e := col.
			Find(query).
			Skip(calculateOffset(pagination)).
			Limit(pagination.Size).
			All(&page.Planets)

		if e != nil {
			return e
		}

		page.Page, e = calculatePage(col, query, pagination, len(page.Planets))
		return e
	})

	return page, err
}

func (pr *planets) Update(planet *models.Planet) errors.Error {
	gomol.Debug("Updating planet on mongodb")

	query := bson.M{
		"_id": planet.ID,
	}

	return pr.execute(func(col *mgo.Collection) error {
		return col.Update(query, bson.M{"$set": planet})
	})
}

func (pr *planets) Delete(id int) errors.Error {
	gomol.Debug("Deleting planet on mongodb")

	query := bson.M{
		"_id": id,
	}

	return pr.execute(func(col *mgo.Collection) error {
		return col.Remove(query)
	})
}
