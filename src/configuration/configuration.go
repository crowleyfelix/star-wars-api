package configuration

import (
	"os"
	"strconv"
)

//Configuration is the application configuration
type Configuration struct {
	Port    int
	MongoDB MongoDB
}

func load() {
	config = new(Configuration)
	config.Port, _ = strconv.Atoi(os.Getenv("PORT"))
	config.MongoDB.Database = os.Getenv("MONGO_DATABASE")
	config.MongoDB.MaxPoolSize, _ = strconv.Atoi(os.Getenv("MONGO_MAX_POOL_SIZE"))
	config.MongoDB.URI = os.Getenv("MONGO_URI")
}

//Get returns application configuration loaded
func Get() *Configuration {
	return config
}
