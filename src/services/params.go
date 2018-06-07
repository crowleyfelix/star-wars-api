package services

import (
	"github.com/crowleyfelix/star-wars-api/src/mongodb"
)

type Pagination struct {
	mongodb.Pagination
}

type PlanetSearchParams struct {
	mongodb.PlanetSearchQuery
}
