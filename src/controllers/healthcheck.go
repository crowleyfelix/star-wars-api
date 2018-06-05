package controllers

import (
	"net/http"

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
	h.context.JSON(http.StatusOK, Response{"I'm alive"})
}
