package main

import (
	"fmt"

	"github.com/crowleyfelix/star-wars-api/server/configuration"
	"github.com/crowleyfelix/star-wars-api/server/controllers"
	"github.com/gin-gonic/gin"
)

type Server interface {
	Start()
}

type server struct {
	engine *gin.Engine
}

//NewServer creates a web server
func NewServer() Server {
	return &server{
		engine: gin.Default(),
	}
}

func (s *server) Start() {
	setUp()
	s.registerRoutes()

	port := fmt.Sprintf(":%d", configuration.Get().Port)

	err := s.engine.Run(port)

	if err != nil {
		panic(err.Error())
	}
}

func (s *server) registerRoutes() {
	controllers.RegisterRoutes(s.engine.Group("/"))
}
