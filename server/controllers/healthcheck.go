package controllers

import (
	"net/http"

	"gopkg.in/mgo.v2"

	"github.com/crowleyfelix/star-wars-api/api/clients/swapi"
	"github.com/crowleyfelix/star-wars-api/api/database/mongodb"
	"github.com/gin-gonic/gin"
)

type healthCheck struct {
	baseController
}

//HealthCheck handles with health check request
func HealthCheck(c *gin.Context) {

	invokeMethod(c, &healthCheck{baseController{c}})
}

func (h *healthCheck) Get() {
	var dependencies []Dependence

	checkers := []checker{
		h.checkSwapi,
		h.checkMongoDB,
	}

	depc := make(chan Dependence, len(checkers))

	go func() {
		for i := range checkers {
			go checkers[i](depc)
		}
	}()

	status := http.StatusOK
	for i := 0; i < len(checkers); i++ {
		dependencies = append(dependencies, <-depc)

		if dependencies[i].Error != nil {
			status = http.StatusInternalServerError
		}
	}

	h.context.JSON(status, dependencies)
}

func (h *healthCheck) checkSwapi(depc chan Dependence) {
	dep := Dependence{Name: "Swapi"}

	if _, err := swapi.New().Endpoints(); err != nil {
		e := err.Error()
		dep.Error = &e
	}

	depc <- dep
}

func (h *healthCheck) checkMongoDB(depc chan Dependence) {
	dep := Dependence{Name: "MongoDB"}
	defer func() {
		depc <- dep
	}()

	var (
		session *mgo.Session
		err     error
	)

	if session, err = mongodb.Pool().Session(); err != nil {
		e := err.Error()
		dep.Error = &e
		return
	}
	if err = session.Ping(); err != nil {
		e := err.Error()
		dep.Error = &e
		return
	}

	mongodb.Pool().Release(session)
}

type checker func(depc chan Dependence)
