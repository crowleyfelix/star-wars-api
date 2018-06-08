package collections

import (
	"fmt"

	"github.com/crowleyfelix/star-wars-api/api/errors"

	"github.com/crowleyfelix/star-wars-api/api/configuration"
	"github.com/crowleyfelix/star-wars-api/api/database/mongodb/models"
	"github.com/golang/glog"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//Planets exposes methods of planet CRUD operations
type Planets interface {
	Insert(*models.Planet) errors.Error
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
	config := configuration.Get().MongoDB
	return &planets{
		collection{
			DataBase:   config.Database,
			Collection: "planets",
			CounterID:  "planet_id",
		},
	}
}

func (pr *planets) Insert(planet *models.Planet) errors.Error {
	glog.Infof("Inserting planet %#v on database", planet)

	return pr.execute(func(col *mgo.Collection) error {
		var err errors.Error
		planet.ID, err = pr.calculateNextID(col.Database)

		if err != nil {
			glog.Errorf("Failed on calculating next id: %s", err.Error())
			return err
		}

		return col.Insert(planet)
	})
}

func (pr *planets) FindByID(id int) (*models.Planet, errors.Error) {
	glog.Infof("Finding planet %d on database", id)

	query := &PlanetSearchQuery{
		ID: &id,
	}
	pagination := &Pagination{
		Page: 1,
		Size: 1,
	}

	page, err := pr.Find(query, pagination)

	if err != nil {
		glog.Errorf("Failed on finding planet id %d on database: %s", id, err.Error())
		return nil, err
	}

	if page.Size == 0 {
		glog.Errorf("Planet id %d was not found on database", id)
		return nil, errors.NewNotFound(fmt.Sprintf("Planet id %d was not found", id))
	}

	return &page.Planets[0], err
}

func (pr *planets) Find(query *PlanetSearchQuery, pagination *Pagination) (*models.PlanetPage, errors.Error) {

	var (
		err  errors.Error
		page = new(models.PlanetPage)
	)

	err = pr.execute(func(col *mgo.Collection) error {
		e := col.
			Find(query).
			Skip(pr.calculateOffset(pagination)).
			Limit(pagination.Size).
			All(&page.Planets)

		if e != nil {
			return e
		}

		page.Page, e = pr.calculatePage(col, query, pagination, len(page.Planets))
		return e
	})

	return page, err
}

func (pr *planets) Update(planet *models.Planet) errors.Error {
	glog.Infof("Updating planet %d on database", planet.ID)

	query := bson.M{
		"_id": planet.ID,
	}

	return pr.execute(func(col *mgo.Collection) error {
		return col.Update(query, bson.M{"$set": planet})
	})
}

func (pr *planets) Delete(id int) errors.Error {
	glog.Infof("Deleting planet %d on database", id)

	query := bson.M{
		"_id": id,
	}

	return pr.execute(func(col *mgo.Collection) error {
		return col.Remove(query)
	})
}
