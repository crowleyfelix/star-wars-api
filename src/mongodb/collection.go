package mongodb

import mgo "gopkg.in/mgo.v2"

type collection struct {
	DataBase   string
	Collection string
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
