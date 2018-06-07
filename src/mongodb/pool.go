package mongodb

import (
	"github.com/crowleyfelix/star-wars-api/src/configuration"
	"github.com/golang/glog"
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
		glog.Info("Creating new session")

		uri := configuration.Get().MongoDB.URI
		p.session, err = mgo.Dial(uri)

		if err != nil {
			glog.Error("Failed on creating session")

			return nil, err
		}
	}

	p.active <- 1
	glog.Info("%d active sessions", len(p.active))

	return p.session.Copy(), nil
}

func (p *pool) Release(session *mgo.Session) {
	glog.Info("Releasing session")

	<-p.active

	glog.Info("%d active sessions", len(p.active))
	session.Close()
}
