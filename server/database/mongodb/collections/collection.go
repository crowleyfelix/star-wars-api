package collections

import (
	"github.com/crowleyfelix/star-wars-api/api/errors"

	"github.com/aphistic/gomol"

	"github.com/crowleyfelix/star-wars-api/api/database/mongodb"
	"github.com/crowleyfelix/star-wars-api/api/database/mongodb/models"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type collection struct {
	DataBase   string
	Collection string
	CounterID  string
}

func (c *collection) execute(operation func(*mgo.Collection) error) errors.Error {
	session, err := mongodb.Pool().Session()

	if err != nil {
		return c.parseError(err)
	}

	defer mongodb.Pool().Release(session)

	col := session.DB(c.DataBase).C(c.Collection)
	return c.parseError(operation(col))
}

func (c *collection) calculateNextID(db *mgo.Database) (int, errors.Error) {
	var counter models.Counter

	change := mgo.Change{
		Update: bson.M{
			"$inc": bson.M{
				"sequence_value": 1,
			},
		},
		ReturnNew: true,
	}

	gomol.Debugf("Updating counter of %s: %#v", c.CounterID, change)

	info, err := db.C("counters").
		Find(bson.M{
			"_id": c.CounterID,
		}).
		Apply(change, &counter)

	if err != nil {
		return 0, c.parseError(err)
	}

	if info.Matched == 0 {
		return 0, errors.NewInternalServer("failed calculating id sequence")
	}

	return counter.SequenceValue, nil
}

func (c *collection) calculatePage(collection *mgo.Collection, query interface{}, pagination *Pagination, totalItems int) (*models.Page, errors.Error) {
	count, err := collection.Find(query).Count()

	if err != nil {
		return nil, c.parseError(err)
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

	return page, nil
}

func (c *collection) calculateOffset(pagination *Pagination) int {
	return (pagination.Page - 1) * pagination.Size
}

func (c *collection) parseError(err error) errors.Error {

	if err == nil {
		return nil
	}

	switch err {
	case mgo.ErrNotFound:
		return errors.NewNotFound(err.Error())
	}

	switch err.(type) {
	case *mgo.QueryError:
	case *mgo.LastError:
		return errors.NewUnprocessableEntity(err.Error())
	}

	return errors.NewInternalServer(err.Error())
}
