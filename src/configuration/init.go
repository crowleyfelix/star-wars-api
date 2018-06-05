package configuration

import (
	"os"
	"strconv"
)

var (
	config *Configuration
)

func init() {
	load()
}

func load() {
	config = new(Configuration)
	config.Port, _ = strconv.Atoi(os.Getenv("PORT"))
}

//Get returns application configuration loaded
func Get() *Configuration {
	return config
}
