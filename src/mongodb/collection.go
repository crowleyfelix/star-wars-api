package mongodb

import (
	"errors"

	"github.com/crowleyfelix/star-wars-api/src/mongodb/models"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type collection struct {
	DataBase   string
	Collection string
	CounterID  string
}

func (c *collection) execute(operation func(*mgo.Collection) error) error {
	session, err := Pool.Session()

	if err != nil {
		return err
	}

	defer Pool.Release(session)

	col := session.DB(c.DataBase).C(c.Collection)
	return operation(col)
}

func (c *collection) calculateNextID(db *mgo.Database) (int, error) {

	var counter models.Counter

	info, err := db.C("counters").
		Find(bson.M{
			"_id": c.CounterID,
		}).
		Apply(mgo.Change{
			Update: bson.M{
				"$inc": bson.M{
					"sequence_value": 1,
				},
			},
			ReturnNew: true,
		}, &counter)

	if err != nil {
		return 0, err
	}

	if info.Matched == 0 {
		return 0, errors.New("failed calculating id sequence")
	}

	return counter.SequenceValue, nil
}
