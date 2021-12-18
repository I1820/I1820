package request

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

// Project request payload.
type Project struct {
	Name        string            `json:"name"`        // project name
	Owner       string            `json:"owner"`       // project owner email address
	Envs        map[string]string `json:"envs"`        // project environment variables
	Description string            `json:"description"` // project description
	Perimeter   []Location        `json:"perimeter"`   // project operational perimeter
}

func (p Project) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.Name, validation.Required),
		validation.Field(&p.Owner, validation.Required, is.Email),
	)
}

type ProjectName struct {
	Name string `json:"name"`
}
