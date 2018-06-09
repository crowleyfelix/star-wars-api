package main

import (
	"fmt"

	"github.com/aphistic/gomol"
	gc "github.com/aphistic/gomol-console"
	"github.com/crowleyfelix/star-wars-api/server/configuration"
	"github.com/crowleyfelix/star-wars-api/server/controllers"
	"github.com/crowleyfelix/star-wars-api/server/middlewares"
	"github.com/gin-gonic/gin"
)

func main() {
	start()
}

func start() {
	setUp()

	engine := gin.Default()
	registerRoutes(engine)

	port := fmt.Sprintf(":%d", configuration.Get().Port)

	err := engine.Run(port)

	if err != nil {
		panic(err.Error())
	}
}

func setUp() {
	if config := configuration.Get(); !config.IsDebug {
		gin.SetMode("release")
	}

	consoleCfg := gc.NewConsoleLoggerConfig()
	consoleLogger, _ := gc.NewConsoleLogger(consoleCfg)
	consoleLogger.SetTemplate(gc.NewTemplateFull())

	gomol.AddLogger(consoleLogger)
	gomol.InitLoggers()

	ch := make(chan error)

	go func() {
		for err := range ch {
			fmt.Printf("[Internal Error] %s\n", err.Error())
		}
	}()

	gomol.SetErrorChan(ch)
}

func registerRoutes(engine *gin.Engine) {
	root := engine.Group("/api")
	middlewares.Register(root)
	controllers.Register(root)
}
