package mongodb

import (
	"github.com/crowleyfelix/star-wars-api/src/configuration"
)

var (
	Pool SessionManager
)

func init() {
	config := configuration.Get().MongoDB
	Pool = &pool{
		active: make(chan int, config.MaxPoolSize),
	}
}
