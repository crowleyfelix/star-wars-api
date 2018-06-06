package models

//Planet represents a star wars planet
type Planet struct {
	ID      int    `bson:"id"`
	Name    string `bson:"name"`
	Climate string `bson:"climate"`
	Terrain string `bson:"terrain"`
}
