package configuration

import (
	"github.com/aphistic/gomol"
	"github.com/crgimenes/goconfig"
)

//Configuration is the application configuration
type Configuration struct {
	Port    int  `cfgDefault:"8888"`
	IsDebug bool `cfg:"debug" cfgDefault:"true"`
	MongoDB MongoDB
}

func load() {
	gomol.Debug("Loading configurations of environment")

	config = new(Configuration)

	err := goconfig.Parse(config)
	if err != nil {
		panic(err.Error())
	}
}

//Get returns application configuration loaded
func Get() *Configuration {
	if config == nil {
		load()
	}

	return config
}
