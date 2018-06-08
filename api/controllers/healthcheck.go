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
	var dependencies []dependencie

	checkers := []checker{
		h.checkSwapi,
		h.checkMongoDB,
	}

	depc := make(chan dependencie, len(checkers))

	go func() {
		for i := range checkers {
			go checkers[i](depc)
		}
	}()

	for i := 0; i < len(checkers); i++ {
		dependencies = append(dependencies, <-depc)
	}

	h.context.JSON(http.StatusOK, dependencies)
}

func (h *healthCheck) checkSwapi(depc chan dependencie) {
	dep := dependencie{Name: "Swapi"}

	if _, err := swapi.New().Endpoints(); err != nil {
		*dep.Error = err.Error()
	}

	depc <- dep
}

func (h *healthCheck) checkMongoDB(depc chan dependencie) {
	dep := dependencie{Name: "MongoDB"}

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

type dependencie struct {
	Name  string  `json:"name"`
	Error *string `json:"error"`
}

type checker func(depc chan dependencie)
