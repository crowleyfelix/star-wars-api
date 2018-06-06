package mongodb

import (
	"github.com/crowleyfelix/star-wars-api/src/configuration"
	"github.com/crowleyfelix/star-wars-api/src/mongodb/models"
	"gopkg.in/mgo.v2"
)

type PlanetCollection interface {
	Find(id int) (*models.Planet, error)
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
		},
	}
}

func (pr *planetCollection) Find(id int) (*models.Planet, error) {
	query := Document{
		"id": id,
	}

	var planet models.Planet

	err := pr.execute(func(col *mgo.Collection) error {
		return col.Find(query).One(&planet)
	})

	return &planet, err
}

func (pr *planetCollection) Create(planet *models.Planet) error {

	return pr.execute(func(col *mgo.Collection) error {
		return col.Insert(planet)
	})
}
