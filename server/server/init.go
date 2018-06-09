package server

import (
	"fmt"

	"github.com/aphistic/gomol"
	gc "github.com/aphistic/gomol-console"
)

func setUp() {
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
