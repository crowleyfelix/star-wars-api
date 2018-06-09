package collections

type Pagination struct {
	Page int
	Size int
}

type PlanetSearchQuery struct {
	ID      *int    `bson:"_id,omitempty"`
	Name    *string `bson:"name,omitempty"`
	Climate *string `bson:"climate,omitempty"`
	Terrain *string `bson:"terrain,omitempty"`
}
