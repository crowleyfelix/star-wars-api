package models

type Page struct {
	Previous *int
	Current  int
	Next     *int
	Size     int
	MaxSize  int
}
