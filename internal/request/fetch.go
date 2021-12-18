package request

import validation "github.com/go-ozzo/ozzo-validation/v4"

// Fetch is a data fetching request.
type Fetch struct {
	ThingIDs []string `json:"thing_ids" validate:"required"`
	Since    int64    `json:"since"`
	Until    int64    `json:"until"`
	Limit    int64    `json:"limit"`
	Offset   int64    `json:"offset"`
}

func (f Fetch) Validate() error {
	return validation.ValidateStruct(&f,
		validation.Field(&f.ThingIDs, validation.Required, validation.Length(1, 0)),
	)
}
