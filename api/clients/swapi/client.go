package swapi

import (
	"errors"
	"fmt"

	"encoding/json"
	"sync"

	"github.com/golang/glog"
	"github.com/levigross/grequests"
)

//Client exposes swapi client methods
type Client interface {
	Endpoints() (map[string]interface{}, error)
	PlanetFilms(name string) ([]Film, error)
}

type client struct{}

//New returns a new Swapi Client
func New() Client {
	return new(client)
}

func (s *client) Endpoints() (map[string]interface{}, error) {
	endpoints := make(map[string]interface{})

	resp, err := grequests.Get(swapiURL, nil)

	if err != nil {
		glog.Errorf("Failed on requesting to swapi")

		return nil, err
	}

	if err = resp.JSON(&endpoints); err != nil {
		glog.Errorf("Failed on deserializing swapi response")
		return nil, err
	}

	return endpoints, err
}

func (s *client) PlanetFilms(name string) ([]Film, error) {
	glog.Infof("Search for planet %s on swapi", name)

	var page Page
	var planets []Planet

	resp, err := grequests.Get(
		fmt.Sprintf("%s/%s?search=%s", swapiURL, swapiPlanetsEndpoint, name),
		nil,
	)

	if err != nil {
		glog.Errorf("Failed on requesting to swapi")

		return nil, err
	}

	if err = resp.JSON(&page); err != nil {
		glog.Errorf("Failed on deserializing swapi response")
		return nil, err
	}

	if err = json.Unmarshal(page.Results, &planets); err != nil {
		glog.Errorf("Failed on deserializing planets result")
		return nil, err
	}

	if len(planets) != 1 {
		glog.Errorf("Planet %s was not found", name)
		return nil, errors.New("not found")
	}

	return s.films(planets[0].Films)
}

func (s *client) films(urls []string) ([]Film, error) {
	group := new(sync.WaitGroup)
	mutex := new(sync.Mutex)

	var err error
	var films []Film

	for i := range urls {
		group.Add(1)

		go func(url string) {
			defer mutex.Unlock()
			defer group.Done()

			film, e := s.film(url)

			if e != nil {
				err = e
			}

			mutex.Lock()
			films = append(films, *film)

		}(urls[i])
	}

	group.Wait()
	if err != nil {
		return nil, err
	}

	return films, nil
}

func (s *client) film(url string) (*Film, error) {
	glog.Infof("Search for film in %s url", url)

	var film Film

	resp, err := grequests.Get(url, nil)

	if err != nil {
		glog.Errorf("Failed on requesting to swapi")
		return nil, err
	}

	if !resp.Ok {
		glog.Errorf("Film %s was not found", url)
		return nil, errors.New("not found")
	}

	if err = resp.JSON(&film); err != nil {
		glog.Errorf("Failed on deserializing swapi response")
		return nil, err
	}

	return &film, nil
}
