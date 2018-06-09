package middlewares

import (
	"github.com/gin-gonic/gin"
)

//Register register API middlewares
func Register(group gin.IRouter) {
	group.Use(Logger)
}
