package controllers

import (
	"github.com/gin-gonic/gin"
)

//Register register API routes
func Register(group gin.IRouter) {
	group.GET("/health-check", HealthCheck)

	planets := group.Group("/planets")
	planets.GET("/", Planets)
	planets.POST("/", Planet)
	planets.GET("/:id", Planet)
	planets.DELETE("/:id", Planet)
}
