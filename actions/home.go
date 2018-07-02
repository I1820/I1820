package actions

import (
	"net/http"

	"github.com/gobuffalo/buffalo"
)

// AboutHandler is a default handler to serve up
// memory of 18.20.
func AboutHandler(c buffalo.Context) error {
	return c.Render(http.StatusOK, r.JSON("18.20 is leaving us"))
}
