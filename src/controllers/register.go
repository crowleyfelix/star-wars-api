package controllers

import (
	"github.com/gin-gonic/gin"
)

//RegisterRoutes register API routes
func RegisterRoutes(group *gin.RouterGroup) {
	group.GET("/health-check", HealthCheck)
}
