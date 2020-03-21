package request

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

// Thing is a request payload for creating thing
type Thing struct {
	Name     string `json:"name"`
	Model    string `json:"model"`
	Location struct {
		Latitude  float64 `json:"lat"`
		Longitude float64 `json:"long"`
	} `json:"location"`
}

func (t Thing) Validate() error {
	return validation.ValidateStruct(&t,
		validation.Field(&t.Name, validation.Required),
		validation.Field(&t.Model, is.Alphanumeric),
		validation.Field(&t.Location.Latitude, is.Latitude),
		validation.Field(&t.Location.Longitude, is.Longitude),
	)
}
