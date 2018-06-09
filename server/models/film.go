package models

import "github.com/crowleyfelix/star-wars-api/server/clients/swapi"

//Film represents an film application model
type Film struct {
	Title     string `json:"title"`
	EpisodeID int64  `json:"episodeId"`
}

//From maps swapi film model to application model
func (f *Film) From(raw *swapi.Film) {
	f.EpisodeID = raw.EpisodeID
	f.Title = raw.Title
}

//Films represents an film collection application model
type Films []Film

//From maps swapi film collection model to application model
func (f *Films) From(raw []swapi.Film) {
	if f == nil {
		temp := make(Films, 0)
		f = &temp
	}

	for _, item := range raw {

		var film Film
		film.From(&item)

		*f = append(*f, film)
	}
}
