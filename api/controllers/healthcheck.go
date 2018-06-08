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
	var dependencies []Dependencie

	checkers := []checker{
		h.checkSwapi,
		h.checkMongoDB,
	}

	depc := make(chan Dependencie, len(checkers))

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

func (h *healthCheck) checkSwapi(depc chan Dependencie) {
	dep := Dependencie{Name: "Swapi"}

	if _, err := swapi.New().Endpoints(); err != nil {
		*dep.Error = err.Error()
	}

	depc <- dep
}

func (h *healthCheck) checkMongoDB(depc chan Dependencie) {
	dep := Dependencie{Name: "MongoDB"}

	var (
		session *mgo.Session
		err     error
	)

	if session, err = mongodb.Pool().Session(); err != nil {
		*dep.Error = err.Error()
	}

	if err = session.Ping(); err != nil {
		*dep.Error = err.Error()
	}

	depc <- dep
}

type checker func(depc chan Dependencie)
