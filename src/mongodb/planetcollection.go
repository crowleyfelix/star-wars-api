package mongodb

import (
	"errors"

	"github.com/crowleyfelix/star-wars-api/src/configuration"
	"github.com/crowleyfelix/star-wars-api/src/mongodb/models"
	"github.com/golang/glog"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//PlanetCollection exposes methods of planet CRUD operations
type PlanetCollection interface {
	Insert(*models.Planet) error
	Find(*PlanetSearchQuery, *Pagination) (*models.PlanetPage, error)
	FindByID(int) (*models.Planet, error)
	Update(*models.Planet) error
	Delete(int) error
}

type planetCollection struct {
	collection
}

//NewPlanetCollection returns new instance of planet collection
func NewPlanetCollection() PlanetCollection {
	config := configuration.Get().MongoDB
	return &planetCollection{
		collection{
			DataBase:   config.Database,
			Collection: "planets",
			CounterID:  "planet_id",
		},
	}
}

func (pr *planetCollection) Insert(planet *models.Planet) error {
	return pr.execute(func(col *mgo.Collection) error {
		var err error
		planet.ID, err = pr.calculateNextID(col.Database)

		if err != nil {
			glog.Errorf("Failed calculating next id: %s", err.Error())
			return err
		}

		return col.Insert(planet)
	})
}

func (pr *planetCollection) FindByID(id int) (*models.Planet, error) {
	query := &PlanetSearchQuery{
		ID: &id,
	}
	pagination := &Pagination{
		Page: 1,
		Size: 1,
	}

	page, err := pr.Find(query, pagination)

	if err != nil {
		return nil, err
	}

	if page.Size == 0 {
		return nil, errors.New("not found")
	}

	return &page.Planets[0], err
}

func (pr *planetCollection) Find(query *PlanetSearchQuery, pagination *Pagination) (*models.PlanetPage, error) {

	var (
		err  error
		page = new(models.PlanetPage)
	)

	err = pr.execute(func(col *mgo.Collection) error {
		err = col.
			Find(query).
			Skip(pr.calculateOffset(pagination)).
			Limit(pagination.Size).
			All(&page.Planets)

		if err != nil {
			return err
		}

		page.Page, err = pr.calculatePage(col, query, pagination, len(page.Planets))
		return err
	})

	return page, err
}

func (pr *planetCollection) Update(planet *models.Planet) error {
	query := bson.M{
		"_id": planet.ID,
	}

	return pr.execute(func(col *mgo.Collection) error {
		return col.Update(query, bson.M{"$set": planet})
	})
}

func (pr *planetCollection) Delete(id int) error {
	query := bson.M{
		"_id": id,
	}

	return pr.execute(func(col *mgo.Collection) error {
		return col.Remove(query)
	})
}
