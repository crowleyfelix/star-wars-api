package mongodb

import (
	"github.com/crowleyfelix/star-wars-api/src/configuration"
	"gopkg.in/mgo.v2"
)

type SessionManager interface {
	Session() (*mgo.Session, error)
	Release(*mgo.Session)
}

type pool struct {
	session *mgo.Session
	active  chan int
}

func (p *pool) Session() (*mgo.Session, error) {
	var err error

	if p.session == nil {
		uri := configuration.Get().MongoDB.URI
		p.session, err = mgo.Dial(uri)

		if err != nil {
			return nil, err
		}
	}

	p.active <- 1
	return p.session.Copy(), nil
}

func (p *pool) Release(session *mgo.Session) {
	<-p.active
	session.Close()
}
