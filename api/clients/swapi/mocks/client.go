// Code generated by mockery v1.0.0
package mocks

import errors "github.com/crowleyfelix/star-wars-api/api/errors"
import mock "github.com/stretchr/testify/mock"
import swapi "github.com/crowleyfelix/star-wars-api/api/clients/swapi"

// Client is an autogenerated mock type for the Client type
type Client struct {
	mock.Mock
}

// Endpoints provides a mock function with given fields:
func (_m *Client) Endpoints() (map[string]interface{}, errors.Error) {
	ret := _m.Called()

	var r0 map[string]interface{}
	if rf, ok := ret.Get(0).(func() map[string]interface{}); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]interface{})
		}
	}

	var r1 errors.Error
	if rf, ok := ret.Get(1).(func() errors.Error); ok {
		r1 = rf()
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(errors.Error)
		}
	}

	return r0, r1
}

// PlanetFilms provides a mock function with given fields: name
func (_m *Client) PlanetFilms(name string) ([]swapi.Film, errors.Error) {
	ret := _m.Called(name)

	var r0 []swapi.Film
	if rf, ok := ret.Get(0).(func(string) []swapi.Film); ok {
		r0 = rf(name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]swapi.Film)
		}
	}

	var r1 errors.Error
	if rf, ok := ret.Get(1).(func(string) errors.Error); ok {
		r1 = rf(name)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(errors.Error)
		}
	}

	return r0, r1
}
