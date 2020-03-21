package request

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

const (
	minLatitude = -90.0
	maxLatitude = 90.0

	minLongitude = 0
	maxLongitude = 180
)

// Location contains latitude and longitude for representing a location
type Location struct {
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"long"`
}

func (l Location) Validate() error {
	return validation.ValidateStruct(&l,
		validation.Field(&l.Latitude, validation.Min(minLatitude), validation.Max(maxLatitude)),
		validation.Field(&l.Longitude, validation.Min(minLongitude), validation.Max(maxLongitude)),
	)
}

// Thing is a request payload for creating thing
type Thing struct {
	Name     string `json:"name"`
	Model    string `json:"model"`
	Location `json:"location"`
}

func (t Thing) Validate() error {
	return validation.ValidateStruct(&t,
		validation.Field(&t.Name, validation.Required),
		validation.Field(&t.Model, is.Alphanumeric),
		validation.Field(&t.Location),
	)
}
