package mongodb

import (
	"github.com/crowleyfelix/star-wars-api/src/configuration"
	"github.com/crowleyfelix/star-wars-api/src/mongodb/models"
	"github.com/golang/glog"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type PlanetCollection interface {
	Insert(planet *models.Planet) error
	Find(id int) (*models.Planet, error)
	List(*Pagination) ([]models.Planet, error)
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

func (pr *planetCollection) Find(id int) (*models.Planet, error) {
	query := bson.M{
		"_id": id,
	}

	var planet models.Planet

	err := pr.execute(func(col *mgo.Collection) error {
		return col.Find(query).One(&planet)
	})

	return &planet, err
}

func (pr *planetCollection) List(pagination *Pagination) ([]models.Planet, error) {

	query := bson.M{}

	var planets []models.Planet

	offset := (pagination.Page - 1) * pagination.Size

	err := pr.execute(func(col *mgo.Collection) error {
		return col.
			Find(query).
			Skip(offset).
			Limit(pagination.Size).
			All(&planets)
	})

	return planets, err
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
