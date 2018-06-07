package mongodb

import (
	"errors"

	"github.com/crowleyfelix/star-wars-api/src/configuration"
	"github.com/crowleyfelix/star-wars-api/src/mongodb/models"
	"github.com/golang/glog"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type PlanetCollection interface {
	Insert(planet *models.Planet) error
	Find(query *PlanetSearchQuery) ([]models.Planet, error)
	FindByID(id int) (*models.Planet, error)
	List(*Pagination) (*models.PlanetPage, error)
	Update(planet *models.Planet) error
	Delete(id int) error
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

	results, err := pr.Find(query)

	if err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, errors.New("not found")
	}

	return &results[0], err
}

func (pr *planetCollection) Find(query *PlanetSearchQuery) ([]models.Planet, error) {
	var planets []models.Planet

	err := pr.execute(func(col *mgo.Collection) error {
		return col.Find(query).All(&planets)
	})

	return planets, err
}

func (pr *planetCollection) List(pagination *Pagination) (*models.PlanetPage, error) {

	var (
		err   error
		query = bson.M{}
		page  = new(models.PlanetPage)
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

		page.Page, err = pr.calculatePage(col, pagination, len(page.Planets))
		return err
	})

	return page, err
}

func (pr *planetCollection) Update(planet *models.Planet) error {
	query := bson.M{
		"_id": planet.ID,
	}

	return pr.execute(func(col *mgo.Collection) error {
		return col.Update(query, planet)
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
