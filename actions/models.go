package actions

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"plugin"

	linkapp "github.com/I1820/link/app"
	"github.com/gobuffalo/buffalo"
)

// ModelsResource represetns model plugins route collection
type ModelsResource struct {
	buffalo.Resource
}

// List returns list of loaded models. This function is mapped to the path
// GET /models
func (m ModelsResource) List(c buffalo.Context) error {
	return c.Render(http.StatusOK, r.JSON(linkApp.Models()))
}

// Create creates model based on given go plugin.
// https://golang.org/pkg/plugin/
// This function is mapped to the path
// POST /models
func (m ModelsResource) Create(c buffalo.Context) error {
	f, err := c.File("model.so")
	if err != nil {
		return c.Error(http.StatusBadRequest, err)
	}

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return c.Error(http.StatusInternalServerError, err)
	}

	if err := ioutil.WriteFile(fmt.Sprintf("./upload/%s", f.Filename), b, 0644); err != nil {
		return c.Error(http.StatusInternalServerError, err)
	}

	p, err := plugin.Open(fmt.Sprintf("./upload/%s", f.Filename))
	if err != nil {
		return c.Error(http.StatusInternalServerError, err)
	}

	sym, err := p.Lookup("Model")
	if err != nil {
		return c.Error(http.StatusInternalServerError, err)
	}
	model, ok := sym.(linkapp.Model)
	if !ok {
		return c.Error(http.StatusBadRequest, fmt.Errorf("Model is required"))
	}
	// TODO synchronous issues?
	linkApp.RegisterModel(model)
	return c.Render(http.StatusOK, r.JSON(model.Name()))
}
