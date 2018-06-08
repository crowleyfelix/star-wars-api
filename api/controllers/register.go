package controllers

import (
	"github.com/gin-gonic/gin"
)

//RegisterRoutes register API routes
func RegisterRoutes(group *gin.RouterGroup) {
	group.GET("/health-check", HealthCheck)

	planets := group.Group("/planets")
	planets.GET("/", Planets)
	planets.POST("/", Planet)
	planets.GET("/:id", Planet)
	planets.DELETE("/:id", Planet)
}
