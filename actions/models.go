package actions

import (
	"net/http"

	"github.com/gobuffalo/buffalo"
)

// ModelsResource represetns model plugins route collection
type ModelsResource struct{}

// List returns list of loaded models
func (m ModelsResource) List(c buffalo.Context) error {
	return c.Render(http.StatusOK, r.JSON(linkApp.Models()))
}
