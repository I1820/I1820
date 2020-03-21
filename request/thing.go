package request

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

// Location contains latitude and longitude for representing a location
type Location struct {
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"long"`
}

func (l Location) Validate() error {
	return validation.ValidateStruct(&l,
		validation.Field(&l.Latitude, is.Latitude),
		validation.Field(&l.Longitude, is.Longitude),
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
