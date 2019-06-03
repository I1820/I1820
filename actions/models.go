package actions

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"plugin"

	"github.com/I1820/link/core"
	"github.com/I1820/link/models"
	"github.com/labstack/echo/v4"
)

// ModelsHandler handles models of core application
type ModelsHandler struct {
	linkApp *core.Application
}

// List returns list of loaded models. This function is mapped to the path
// GET /models
func (m ModelsHandler) List(c echo.Context) error {
	return c.JSON(http.StatusOK, m.linkApp.Models())
}

// Create creates model based on given go plugin.
// https://golang.org/pkg/plugin/
// This function is mapped to the path
// POST /models
func (m ModelsHandler) Create(c echo.Context) error {
	fh, err := c.FormFile("model.so")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	f, err := fh.Open()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if err := ioutil.WriteFile(fmt.Sprintf("./upload/%s", fh.Filename), b, 0644); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	p, err := plugin.Open(fmt.Sprintf("./upload/%s", fh.Filename))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	sym, err := p.Lookup("Model")
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	model, ok := sym.(models.Model)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "please upload a valid Model")
	}
	// TODO synchronous issues?
	m.linkApp.RegisterModel(model)
	return c.JSON(http.StatusOK, model.Name())
}
