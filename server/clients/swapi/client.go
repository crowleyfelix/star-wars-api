package swapi

import (
	"fmt"

	"encoding/json"
	"sync"

	"github.com/crowleyfelix/star-wars-api/server/errors"

	"github.com/aphistic/gomol"
	"github.com/levigross/grequests"
)

//Client exposes swapi client methods
type Client interface {
	Endpoints() (map[string]interface{}, errors.Error)
	Planet(name string) (*Planet, errors.Error)
	PlanetFilms(name string) ([]Film, errors.Error)
}

type client struct{}

//New returns a new Swapi Client
func New() Client {
	gomol.Debug("Creating new Swapi Client")
	return new(client)
}

func (s *client) Endpoints() (map[string]interface{}, errors.Error) {
	endpoints := make(map[string]interface{})

	resp, err := grequests.Get(swapiURL, nil)

	if err != nil {
		gomol.Errorf("Failed on requesting to swapi")
		return nil, errors.NewCallout(err.Error())
	}

	if err = resp.JSON(&endpoints); err != nil {
		gomol.Errorf("Failed on deserializing swapi response")
		return nil, errors.NewSerialization(err.Error())
	}

	return endpoints, nil
}

func (s *client) Planet(name string) (*Planet, errors.Error) {
	gomol.Debugf("Searching for planet %s on swapi", name)

	var page Page
	var planets []Planet

	resp, err := grequests.Get(
		fmt.Sprintf("%s/%s?search=%s", swapiURL, swapiPlanetsEndpoint, name),
		nil,
	)

	if err != nil {
		gomol.Errorf("Failed on requesting to swapi")
		return nil, errors.NewCallout(err.Error())
	}

	if err = resp.JSON(&page); err != nil {
		gomol.Errorf("Failed on deserializing swapi response")
		return nil, errors.NewSerialization(err.Error())
	}

	if err = json.Unmarshal(page.Results, &planets); err != nil {
		gomol.Errorf("Failed on deserializing planets result")
		return nil, errors.NewSerialization(err.Error())
	}

	if len(planets) != 1 {
		gomol.Errorf("Planet %s was not found on swapi", name)
		return nil, errors.NewNotFound(fmt.Sprintf("Planet %s was not found", name))
	}

	return &planets[0], nil
}

func (s *client) PlanetFilms(name string) ([]Film, errors.Error) {
	gomol.Debugf("Searching for planet %s on swapi", name)

	planet, err := s.Planet(name)

	if err != nil {
		return nil, err
	}

	return s.films(planet.Films)
}

func (s *client) films(urls []string) ([]Film, errors.Error) {
	group := new(sync.WaitGroup)
	mutex := new(sync.Mutex)

	var err errors.Error
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

func (s *client) film(url string) (*Film, errors.Error) {
	gomol.Debugf("Searching for film in %s", url)

	var film Film

	resp, err := grequests.Get(url, nil)

	if err != nil {
		gomol.Errorf("Failed on requesting to swapi")
		return nil, errors.NewCallout(err.Error())
	}

	if !resp.Ok {
		gomol.Errorf("Film %s responded with status code %d", url, resp.StatusCode)
		return nil, errors.Build(resp.StatusCode)
	}

	if err = resp.JSON(&film); err != nil {
		gomol.Errorf("Failed on deserializing swapi response")
		return nil, errors.NewSerialization(err.Error())
	}

	return &film, nil
}
