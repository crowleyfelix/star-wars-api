package services

import (
	mongodb "github.com/crowleyfelix/star-wars-api/api/database/mongodb/collections"
)

type Pagination struct {
	mongodb.Pagination
}

type PlanetSearchParams struct {
	mongodb.PlanetSearchQuery
}
