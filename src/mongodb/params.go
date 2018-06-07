package mongodb

type Pagination struct {
	Page int
	Size int
}

type PlanetSearchQuery struct {
	ID   *int    `bson:"_id,omitempty"`
	Name *string `bson:"name,omitempty"`
}
