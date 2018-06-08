package mongodb

import (
	"github.com/crowleyfelix/star-wars-api/api/configuration"
	"github.com/crowleyfelix/star-wars-api/api/errors"
	"github.com/golang/glog"
	"gopkg.in/mgo.v2"
)

//SessionPool exposes session manager methods
type SessionPool interface {
	Session() (*mgo.Session, errors.Error)
	Release(*mgo.Session)
}

type pool struct {
	session *mgo.Session
	active  chan int
}

//Pool returns mongodb session manager
func Pool() SessionPool {
	if sessionPool == nil {
		config := configuration.Get().MongoDB
		sessionPool = &pool{
			active: make(chan int, config.MaxPoolSize),
		}
	}
	return sessionPool
}

func (p *pool) Session() (*mgo.Session, errors.Error) {
	var err error

	if p.session == nil {
		glog.Info("Creating new session")

		uri := configuration.Get().MongoDB.URI
		p.session, err = mgo.Dial(uri)

		if err != nil {
			glog.Error("Failed on creating session")

			return nil, errors.NewInternalServer(err.Error())
		}
	}

	p.active <- 1
	glog.Infof("%d active sessions", len(p.active))

	return p.session.Copy(), nil
}

func (p *pool) Release(session *mgo.Session) {
	glog.Info("Releasing session")

	<-p.active

	glog.Infof("%d active sessions", len(p.active))
	session.Close()
}
