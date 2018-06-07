package mongodb

import (
	"github.com/crowleyfelix/star-wars-api/src/configuration"
	"github.com/crowleyfelix/star-wars-api/src/mongodb/models"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type PlanetCollection interface {
	Find(id int) (*models.Planet, error)
	Insert(planet *models.Planet) error
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

func (pr *planetCollection) Find(id int) (*models.Planet, error) {
	query := bson.M{
		"id": id,
	}

	var planet models.Planet

	err := pr.execute(func(col *mgo.Collection) error {
		return col.Find(query).One(&planet)
	})

	return &planet, err
}

func (pr *planetCollection) Insert(planet *models.Planet) error {
	return pr.execute(func(col *mgo.Collection) error {
		var err error
		planet.ID, err = pr.calculateNextID(col.Database)

		if err != nil {
			return err
		}

		return col.Insert(planet)
	})
}

func (pr *planetCollection) Update(planet *models.Planet) error {
	query := bson.M{
		"id": planet.ID,
	}

	return pr.execute(func(col *mgo.Collection) error {
		return col.Update(query, planet)
	})
}
