package mongodb

import (
	"errors"

	"github.com/golang/glog"

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
	session, err := Pool().Session()

	if err != nil {
		return err
	}

	defer Pool().Release(session)

	col := session.DB(c.DataBase).C(c.Collection)
	return operation(col)
}

func (c *collection) calculateNextID(db *mgo.Database) (int, error) {
	var counter models.Counter

	change := mgo.Change{
		Update: bson.M{
			"$inc": bson.M{
				"sequence_value": 1,
			},
		},
		ReturnNew: true,
	}

	glog.Infof("Updating counter of %s: %#v", c.CounterID, change)

	info, err := db.C("counters").
		Find(bson.M{
			"_id": c.CounterID,
		}).
		Apply(change, &counter)

	if err != nil {
		return 0, err
	}

	if info.Matched == 0 {
		return 0, errors.New("failed calculating id sequence")
	}

	return counter.SequenceValue, nil
}

func (c *collection) calculatePage(collection *mgo.Collection, query interface{}, pagination *Pagination, totalItems int) (*models.Page, error) {
	count, err := collection.Find(query).Count()

	if err != nil {
		return nil, err
	}

	page := &models.Page{
		Current: pagination.Page,
		Size:    totalItems,
		MaxSize: pagination.Size,
	}
	if pagination.Page > 1 {
		previous := pagination.Page - 1
		page.Previous = &previous
	}
	if pagination.Page*pagination.Size < count {
		next := pagination.Page + 1
		page.Next = &next
	}

	return page, err
}

func (c *collection) calculateOffset(pagination *Pagination) int {
	return (pagination.Page - 1) * pagination.Size
}
