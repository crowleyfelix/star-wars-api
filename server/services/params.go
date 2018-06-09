package services

import (
	mongodb "github.com/crowleyfelix/star-wars-api/server/database/mongodb/collections"
)

type Pagination struct {
	mongodb.Pagination
}

type PlanetSearchParams struct {
	mongodb.PlanetSearchQuery
}
